package api

import (
	"github.com/filecoin-project/venus-wallet/common"
	"github.com/filecoin-project/venus-wallet/storage/strategy"
	"github.com/filecoin-project/venus-wallet/storage/wallet"
	"github.com/filecoin-project/venus-wallet/wallet_event"
)

type IFullAPI interface {
	common.ICommon
	strategy.ILocalStrategy
	wallet.ILocalWallet
	wallet_event.IWalletEvent
}

type FullAPI struct {
	common.ICommon
	strategy.ILocalStrategy
	wallet.ILocalWallet
	wallet_event.IWalletEvent
}
