package api

import (
	"context"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/storage/wallet"
)

var _ wallet.IWallet = &WalletAuth{}

// wallet API permissions constraints
type WalletAuth struct {
	Internal struct {
		WalletNew    func(context.Context, core.KeyType) (core.Address, error)                                                 `perm:"admin"`
		WalletHas    func(ctx context.Context, address core.Address) (bool, error)                                             `perm:"write"`
		WalletList   func(ctx context.Context) ([]core.Address, error)                                                         `perm:"write"`
		WalletSign   func(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error) `perm:"sign"`
		WalletExport func(ctx context.Context, addr core.Address) (*core.KeyInfo, error)                                       `perm:"admin"`
		WalletImport func(context.Context, *core.KeyInfo) (core.Address, error)                                                `perm:"admin"`
		WalletDelete func(context.Context, core.Address) error                                                                 `perm:"admin"`
	}
}

func (c *WalletAuth) WalletNew(ctx context.Context, keyType core.KeyType) (core.Address, error) {
	return c.Internal.WalletNew(ctx, keyType)
}

func (c *WalletAuth) WalletHas(ctx context.Context, addr core.Address) (bool, error) {
	return c.Internal.WalletHas(ctx, addr)
}

func (c *WalletAuth) WalletList(ctx context.Context) ([]core.Address, error) {
	return c.Internal.WalletList(ctx)
}

func (c *WalletAuth) WalletSign(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error) {
	return c.Internal.WalletSign(ctx, signer, toSign, meta)
}

func (c *WalletAuth) WalletExport(ctx context.Context, a core.Address) (*core.KeyInfo, error) {
	return c.Internal.WalletExport(ctx, a)
}

func (c *WalletAuth) WalletImport(ctx context.Context, ki *core.KeyInfo) (core.Address, error) {
	return c.Internal.WalletImport(ctx, ki)
}

func (c *WalletAuth) WalletDelete(ctx context.Context, addr core.Address) error {
	return c.Internal.WalletDelete(ctx, addr)
}
