package api

import (
	"github.com/ipfs-force-community/venus-wallet/storage"
	"github.com/multiformats/go-multiaddr"
)

type APIEndpoint multiaddr.Multiaddr

type FullAPI struct {
	storage.IWallet
	ICommon
}

type IFullAPI interface {
	storage.IWallet
	ICommon
}
