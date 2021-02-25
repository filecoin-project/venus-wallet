package api

import (
	"context"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/multiformats/go-multiaddr"
)

// remote wallet api
type IWallet interface {
	WalletNew(context.Context, core.KeyType) (core.Address, error)
	WalletHas(ctx context.Context, address core.Address) (bool, error)
	WalletList(ctx context.Context) ([]core.Address, error)
	WalletSign(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error)
	WalletExport(ctx context.Context, addr core.Address) (*core.KeyInfo, error)
	WalletImport(context.Context, *core.KeyInfo) (core.Address, error)
	WalletDelete(context.Context, core.Address) error
}

// rpc api endpoint
type APIEndpoint multiaddr.Multiaddr

type FullAPI struct {
	IWallet
	ICommon
}

type IFullAPI interface {
	IWallet
	ICommon
}
