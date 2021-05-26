package api

import (
	"github.com/filecoin-project/venus-wallet/common"
	"github.com/filecoin-project/venus-wallet/storage/strategy"
	"github.com/filecoin-project/venus-wallet/storage/wallet"
	"github.com/filecoin-project/venus-wallet/wallet_event"
	"github.com/multiformats/go-multiaddr"
)

// rpc api endpoint
type APIEndpoint multiaddr.Multiaddr

type FullAPI struct {
	strategy.ILocalStrategy
	wallet.ILocalWallet
	common.ICommon
	wallet_event.IWalletEventAPI
}

type IFullAPI interface {
	strategy.ILocalStrategy
	wallet.ILocalWallet
	common.ICommon
	wallet_event.IWalletEventAPI
}
