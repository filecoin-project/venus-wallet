package api

import (
	"github.com/filecoin-project/venus-wallet/common"
	"github.com/filecoin-project/venus-wallet/storage/strategy"
	"github.com/filecoin-project/venus-wallet/storage/wallet"
	"github.com/multiformats/go-multiaddr"
)

// rpc api endpoint
type APIEndpoint multiaddr.Multiaddr

type FullAPI struct {
	strategy.ILocalStrategy
	wallet.ILocalWallet
	common.ICommon
}

type IFullAPI interface {
	strategy.ILocalStrategy
	wallet.ILocalWallet
	common.ICommon
}
