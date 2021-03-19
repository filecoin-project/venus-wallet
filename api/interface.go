package api

import (
	"github.com/ipfs-force-community/venus-wallet/common"
	"github.com/ipfs-force-community/venus-wallet/storage/strategy"
	"github.com/ipfs-force-community/venus-wallet/storage/wallet"
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
