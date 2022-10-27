package wallet_event

import (
	"context"
	"sync"

	"github.com/filecoin-project/venus/venus-shared/api/gateway/v1"

	"github.com/asaskevich/EventBus"
	"github.com/filecoin-project/go-address"
	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"

	"github.com/filecoin-project/venus-wallet/config"
	"github.com/ipfs-force-community/venus-gateway/types"
	"github.com/ipfs-force-community/venus-gateway/walletevent"
)

var log = logging.Logger("wallet_event")

type IAPIRegisterHub interface {
	SupportNewAccount(ctx context.Context, account string) error
	AddNewAddress(ctx context.Context, newAddrs []address.Address) error
	RemoveAddress(ctx context.Context, newAddrs []address.Address) error
}

type APIRegisterHub struct {
	weClient map[string]*walletevent.WalletEventClient
	bus      EventBus.Bus
	lk       sync.Mutex
}

func NewAPIRegisterHub(lc fx.Lifecycle, signer types.IWalletHandler, bus EventBus.Bus, cfg *config.APIRegisterHubConfig) (*APIRegisterHub, error) {
	apiRegister := &APIRegisterHub{
		weClient: make(map[string]*walletevent.WalletEventClient),
		bus:      bus,
		lk:       sync.Mutex{},
	}

	if len(cfg.RegisterAPI) == 0 {
		log.Warnf("api hub: urls: %v, token: %s, support account: %v", cfg.RegisterAPI, cfg.Token, cfg.SupportAccounts)
	} else {
		log.Infof("api hub: urls: %v, token: %s, support account: %v", cfg.RegisterAPI, cfg.Token, cfg.SupportAccounts)
	}

	for _, apiHub := range cfg.RegisterAPI {
		ctx, cancel := context.WithCancel(context.Background())
		walletEventClient, closer, err := gateway.DialIGatewayRPC(ctx, apiHub, cfg.Token, nil)
		if err != nil {
			// todo return or continue. allow failed client
			log.Errorf("connect to api hub %s failed %v", apiHub, err)
			cancel()
			return nil, err
		}
		mLog := log.With("api hub", apiHub)
		walletEvent := walletevent.NewWalletEventClient(ctx, signer, walletEventClient, mLog, cfg.SupportAccounts)
		go walletEvent.ListenWalletRequest(ctx)
		apiRegister.lk.Lock()
		apiRegister.weClient[apiHub] = walletEvent
		apiRegister.lk.Unlock()
		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				closer()
				cancel()
				return nil
			},
		})
	}

	_ = bus.Subscribe("wallet:add_address", func(addr address.Address) {
		log.Infof("wallet add address %s", addr)
		err := apiRegister.AddNewAddress(context.TODO(), []address.Address{addr})
		if err != nil {
			log.Errorf("cannot add address %s, %s", addr, err.Error())
		}
	})

	_ = bus.Subscribe("wallet:remove_address", func(addr address.Address) {
		log.Infof("wallet remove address %s", addr)
		err := apiRegister.RemoveAddress(context.TODO(), []address.Address{addr})
		if err != nil {
			log.Errorf("cannot remove address %s", addr)
		}
	})
	return apiRegister, nil
}

func (h *APIRegisterHub) SupportNewAccount(ctx context.Context, supportAccount string) error {
	h.lk.Lock()
	defer h.lk.Unlock()
	for _, c := range h.weClient {
		err := c.SupportAccount(ctx, supportAccount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *APIRegisterHub) AddNewAddress(ctx context.Context, newAddrs []address.Address) error {
	h.lk.Lock()
	defer h.lk.Unlock()
	for _, c := range h.weClient {
		err := c.AddNewAddress(ctx, newAddrs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *APIRegisterHub) RemoveAddress(ctx context.Context, newAddrs []address.Address) error {
	h.lk.Lock()
	defer h.lk.Unlock()
	for _, c := range h.weClient {
		err := c.RemoveAddress(ctx, newAddrs)
		if err != nil {
			return err
		}
	}
	return nil
}
