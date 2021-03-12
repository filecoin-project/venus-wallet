package api

import (
	"github.com/ipfs-force-community/venus-wallet/common"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"github.com/ipfs-force-community/venus-wallet/storage/wallet"
	"github.com/multiformats/go-multiaddr"
)

// rpc api endpoint
type APIEndpoint multiaddr.Multiaddr

type FullAPI struct {
	wallet.IWallet
	storage.IWalletLock
	common.ICommon
	storage.StrategyStore
}

type IFullAPI interface {
	wallet.IWallet
	storage.IWalletLock
	common.ICommon
	storage.StrategyStore
}
