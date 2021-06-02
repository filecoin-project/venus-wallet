package wallet_event

import (
	"context"
	"encoding/json"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-wallet/config"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/storage/wallet"
	"github.com/google/uuid"
	"github.com/ipfs-force-community/venus-gateway/types"
	"github.com/ipfs-force-community/venus-gateway/walletevent"
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
	"sync"
	"time"
)

var log = logging.Logger("wallet_event")

type IAPIRegisterHub interface {
	SupportNewAccount(ctx context.Context, account string) error
	AddNewAddress(ctx context.Context, newAddrs []address.Address) error
}

type IWalletProcess interface {
	WalletList(ctx context.Context) ([]core.Address, error)
	WalletSign(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error)
}

type APIRegisterHub struct {
	registerClient map[string]*WalletEvent
	lk             sync.Mutex
}

func NewAPIRegisterHub(lc fx.Lifecycle, process wallet.ILocalWallet, cfg *config.APIRegisterHubConfig) (*APIRegisterHub, error) {
	apiRegister := &APIRegisterHub{
		registerClient: make(map[string]*WalletEvent),
		lk:             sync.Mutex{},
	}

	for _, apiHub := range cfg.RegisterAPI {
		ctx, cancel := context.WithCancel(context.Background())
		walletEventClient, closer, err := NewWalletRegisterClient(ctx, apiHub, cfg.Token)
		if err != nil {
			//todo return or continue. allow failed client
			log.Errorf("connect to api hub %s faile %v", apiHub, err)
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
	}
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
			addrs, err := e.processor.WalletList(ctx)
			if err != nil {
				e.log.Errorf("WalletList error %s", err)
				e.error(ctx, event.Id, err)
				continue
			}

			e.value(ctx, event.Id, addrs)
		case "WalletSign":
			log.Info("receive WalletSign event")
			req := types.WalletSignRequest{}
			err := json.Unmarshal(event.Payload, &req)
			if err != nil {
				e.log.Errorf("unmarshal WalletSignRequest error %s", err)
				e.error(ctx, event.Id, err)
				continue
			}
			log.Info("start WalletSign")
			sig, err := e.processor.WalletSign(ctx, req.Signer, req.ToSign, req.Meta)
			if err != nil {
				e.log.Errorf("WalletSign error %s", err)
				e.error(ctx, event.Id, err)
				continue
			}
			log.Info("end WalletSign")
			e.value(ctx, event.Id, sig)
			log.Info("end WalletSign response")

		default:
			e.log.Errorf("unexpect proof event type %s", event.Method)
		}
	}

	return nil
}

func (e *WalletEvent) value(ctx context.Context, id uuid.UUID, val interface{}) {
	respBytes, err := json.Marshal(val)
	if err != nil {
		e.log.Errorf("marshal address list error %s", err)
		e.client.ResponseWalletEvent(ctx, &types.ResponseEvent{
			Id:      id,
			Payload: nil,
			Error:   err.Error(),
		})
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
