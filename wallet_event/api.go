package wallet_event

import (
	"context"
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-wallet/config"
	"github.com/filecoin-project/venus-wallet/filemgr"
	api "github.com/filecoin-project/venus/venus-shared/api/wallet"
)

type IWalletEvent = api.IWalletEvent

type WalletEventAPI struct {
	cfg            *config.APIRegisterHubConfig
	apiRegisterHub IAPIRegisterHub
	fsr            filemgr.Repo
}

func NewWalletEventAPI(fsr filemgr.Repo, cfg *config.APIRegisterHubConfig, apiRegisterHub IAPIRegisterHub) IWalletEvent {
	return &WalletEventAPI{
		fsr:            fsr,
		cfg:            cfg,
		apiRegisterHub: apiRegisterHub,
	}
}

func (walletEventAPI *WalletEventAPI) AddSupportAccount(ctx context.Context, supportAccount string) error {
	for _, account := range walletEventAPI.cfg.SupportAccounts {
		if account == supportAccount {
			return fmt.Errorf("account %s has exit", supportAccount)
		}
	}

	err := walletEventAPI.apiRegisterHub.SupportNewAccount(ctx, supportAccount)
	if err != nil {
		return err
	}
	return walletEventAPI.fsr.AppendSupportAccount(supportAccount)
}

func (walletEventAPI *WalletEventAPI) AddNewAddress(ctx context.Context, newAddrs []address.Address) error {
	return walletEventAPI.apiRegisterHub.AddNewAddress(ctx, newAddrs)
}
