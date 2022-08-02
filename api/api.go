package api

import (
	api "github.com/filecoin-project/venus/venus-shared/api/wallet"
)

var _ api.IFullAPI = (*FullAPI)(nil)

type FullAPI struct {
	api.ICommon
	api.ILocalStrategy
	api.ILocalWallet
	api.IWalletEvent
}
