package api

import (
	"github.com/filecoin-project/venus-wallet/common"
	"github.com/filecoin-project/venus-wallet/storage/strategy"
	wallet_api "github.com/filecoin-project/venus/venus-shared/api/wallet"
)

type IFullAPI interface {
	common.ICommon
	strategy.ILocalStrategy
	wallet_api.ILocalWallet
	wallet_api.IWalletEvent
}

type FullAPI struct {
	common.ICommon
	strategy.ILocalStrategy
	wallet_api.ILocalWallet
	wallet_api.IWalletEvent
}
