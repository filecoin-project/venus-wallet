package storage

import (
	"context"
	"github.com/ipfs-force-community/venus-wallet/core"
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
