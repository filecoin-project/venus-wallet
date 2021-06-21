package wallet_event

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/filecoin-project/go-address"
	"github.com/google/uuid"
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/venus-wallet/config"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/ipfs-force-community/venus-gateway/types"
	"github.com/ipfs-force-community/venus-gateway/walletevent"
)

type ShimWallet interface {
	WalletList(ctx context.Context) ([]core.Address, error)
	WalletSign(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error)
}

var log = logging.Logger("wallet_event")

type IAPIRegisterHub interface {
	SupportNewAccount(ctx context.Context, account string) error
	AddNewAddress(ctx context.Context, newAddrs []address.Address) error
	RemoveAddress(ctx context.Context, newAddrs []address.Address) error
}

type IWalletProcess interface {
	WalletList(ctx context.Context) ([]core.Address, error)
	WalletSign(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error)
}

type APIRegisterHub struct {
	registerClient map[string]*WalletEvent
	bus            EventBus.Bus
	lk             sync.Mutex
}

func NewAPIRegisterHub(lc fx.Lifecycle, process ShimWallet, bus EventBus.Bus, cfg *config.APIRegisterHubConfig) (*APIRegisterHub, error) {
	apiRegister := &APIRegisterHub{
		registerClient: make(map[string]*WalletEvent),
		bus:            bus,
		lk:             sync.Mutex{},
	}

	for _, apiHub := range cfg.RegisterAPI {
		ctx, cancel := context.WithCancel(context.Background())
		walletEventClient, closer, err := NewWalletRegisterClient(ctx, apiHub, cfg.Token)
		if err != nil {
			//todo return or continue. allow failed client
			log.Errorf("connect to api hub %s failed %v", apiHub, err)
			cancel()
			return nil, err
		}
		mLog := log.With("api hub", apiHub)
		walletEvent := NewWalletEvent(ctx, process, walletEventClient, mLog, cfg)
		go walletEvent.listenWalletRequest(ctx)
		apiRegister.lk.Lock()
		apiRegister.registerClient[apiHub] = walletEvent
		apiRegister.lk.Unlock()
		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				closer()
				cancel()
				return nil
			},
		})
	}

	_ = bus.Subscribe("wallet:add_address", func(addr core.Address) {
		log.Infof("wallet add address %s", addr)
		err := apiRegister.AddNewAddress(context.TODO(), []address.Address{addr})
		if err != nil {
			log.Errorf("cannot add addres %s", addr)
		}
	})

	_ = bus.Subscribe("wallet:remove_address", func(addr core.Address) {
		log.Infof("wallet remove address %s", addr)
		err := apiRegister.RemoveAddress(context.TODO(), []address.Address{addr})
		if err != nil {
			log.Errorf("cannot remove addres %s", addr)
		}
	})
	return apiRegister, nil
}

func (apiRegisterhub *APIRegisterHub) SupportNewAccount(ctx context.Context, supportAccount string) error {
	apiRegisterhub.lk.Lock()
	defer apiRegisterhub.lk.Unlock()
	for _, c := range apiRegisterhub.registerClient {
		err := c.SupportAccount(ctx, supportAccount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (apiRegisterhub *APIRegisterHub) AddNewAddress(ctx context.Context, newAddrs []address.Address) error {
	apiRegisterhub.lk.Lock()
	defer apiRegisterhub.lk.Unlock()
	for _, c := range apiRegisterhub.registerClient {
		err := c.AddNewAddress(ctx, newAddrs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (apiRegisterhub *APIRegisterHub) RemoveAddress(ctx context.Context, newAddrs []address.Address) error {
	apiRegisterhub.lk.Lock()
	defer apiRegisterhub.lk.Unlock()
	for _, c := range apiRegisterhub.registerClient {
		err := c.RemoveAddress(ctx, newAddrs)
		if err != nil {
			return err
		}
	}
	return nil
}

type WalletEvent struct {
	processor IWalletProcess
	client    *WalletRegisterClient
	log       logging.StandardLogger
	channel   uuid.UUID
	cfg       *config.APIRegisterHubConfig
}

func NewWalletEvent(ctx context.Context, process IWalletProcess, client *WalletRegisterClient, log logging.StandardLogger, cfg *config.APIRegisterHubConfig) *WalletEvent {
	return &WalletEvent{processor: process, client: client, log: log, cfg: cfg}
}

func (e *WalletEvent) SupportAccount(ctx context.Context, supportAccount string) error {
	err := e.client.SupportNewAccount(ctx, e.channel, supportAccount)
	if err != nil {
		return err
	}
	return nil
}

func (e *WalletEvent) AddNewAddress(ctx context.Context, newAddrs []address.Address) error {
	return e.client.AddNewAddress(ctx, e.channel, newAddrs)
}

func (e *WalletEvent) RemoveAddress(ctx context.Context, newAddrs []address.Address) error {
	return e.client.RemoveAddress(ctx, e.channel, newAddrs)
}

func (e *WalletEvent) listenWalletRequest(ctx context.Context) {
	for {
		if err := e.listenWalletRequestOnce(ctx); err != nil {
			e.log.Errorf("listen wallet event errored: %s", err)
		} else {
			e.log.Warn("listenWalletRequestOnce quit")
		}
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			e.log.Warnf("not restarting listenWalletRequestOnce: context error: %s", ctx.Err())
			return
		}

		e.log.Info("restarting listenWalletRequestOnce")
	}
}

func (e *WalletEvent) listenWalletRequestOnce(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	policy := &walletevent.WalletRegisterPolicy{
		SupportAccounts: e.cfg.SupportAccounts,
		SignBytes:       core.RandSignBytes,
	}
	log.Infow("", "rand sign byte", core.RandSignBytes)
	walletEventCh, err := e.client.ListenWalletEvent(ctx, policy)
	if err != nil {
		// Retry is handled by caller
		return xerrors.Errorf("listenWalletRequestOnce listenWalletRequestOnce call failed: %w", err)
	}

	for event := range walletEventCh {
		switch event.Method {
		case "InitConnect":
			req := types.ConnectedCompleted{}
			err := json.Unmarshal(event.Payload, &req)
			if err != nil {
				e.log.Errorf("init connect error %s", err)
			}
			e.channel = req.ChannelId
			e.log.Infof("connect to server %v", req.ChannelId)
			//do not response
		case "WalletList":
			go e.walletList(ctx, event.Id)
		case "WalletSign":
			go e.walletSign(ctx, event)
		default:
			e.log.Errorf("unexpect proof event type %s", event.Method)
		}
	}

	return nil
}

func (e *WalletEvent) walletList(ctx context.Context, id uuid.UUID) {
	addrs, err := e.processor.WalletList(ctx)
	if err != nil {
		e.log.Errorf("WalletList error %s", err)
		e.error(ctx, id, err)
		return
	}
	e.value(ctx, id, addrs)
}

func (e *WalletEvent) walletSign(ctx context.Context, event *types.RequestEvent) {
	log.Debug("receive WalletSign event")
	req := types.WalletSignRequest{}
	err := json.Unmarshal(event.Payload, &req)
	if err != nil {
		e.log.Errorf("unmarshal WalletSignRequest error %s", err)
		e.error(ctx, event.Id, err)
		return
	}
	log.Debug("start WalletSign")
	sig, err := e.processor.WalletSign(ctx, req.Signer, req.ToSign, core.MsgMeta{Type: core.MsgType(req.Meta.Type), Extra: req.Meta.Extra})
	if err != nil {
		e.log.Errorf("WalletSign error %s", err)
		e.error(ctx, event.Id, err)
		return
	}
	log.Debug("end WalletSign")
	e.value(ctx, event.Id, sig)
	log.Debug("end WalletSign response")
}

func (e *WalletEvent) value(ctx context.Context, id uuid.UUID, val interface{}) {
	respBytes, err := json.Marshal(val)
	if err != nil {
		e.log.Errorf("marshal address list error %s", err)
		err = e.client.ResponseWalletEvent(ctx, &types.ResponseEvent{
			Id:      id,
			Payload: nil,
			Error:   err.Error(),
		})
		e.log.Errorf("response wallet event error %s", err)
		return
	}
	err = e.client.ResponseWalletEvent(ctx, &types.ResponseEvent{
		Id:      id,
		Payload: respBytes,
		Error:   "",
	})
	if err != nil {
		e.log.Errorf("response error %v", err)
	}
}

func (e *WalletEvent) error(ctx context.Context, id uuid.UUID, err error) {
	err = e.client.ResponseWalletEvent(ctx, &types.ResponseEvent{
		Id:      id,
		Payload: nil,
		Error:   err.Error(),
	})
	if err != nil {
		e.log.Errorf("response error %v", err)
	}
}
