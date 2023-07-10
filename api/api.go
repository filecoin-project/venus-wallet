package api

import (
	"github.com/filecoin-project/venus-wallet/common"
	wallet_api "github.com/filecoin-project/venus/venus-shared/api/wallet"
	"go.uber.org/fx"
)

type IFullAPI interface {
	common.ICommon
	wallet_api.ILocalWallet
	wallet_api.IWalletEvent
}

type FullAPI struct {
	fx.In
	common.ICommon
	wallet_api.ILocalWallet
	wallet_api.IWalletEvent
}
