package api

import (
	"context"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/storage/wallet"
)

var _ wallet.IWallet = &WalletAPIAdapter{}

// wallet API permissions constraints
type WalletAPIAdapter struct {
	Internal struct {
		WalletNew    func(ctx context.Context, kt core.KeyType) (core.Address, error)                                          `perm:"admin"`
		WalletHas    func(ctx context.Context, address core.Address) (bool, error)                                             `perm:"write"`
		WalletList   func(ctx context.Context) ([]core.Address, error)                                                         `perm:"write"`
		WalletSign   func(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error) `perm:"sign"`
		WalletExport func(ctx context.Context, addr core.Address) (*core.KeyInfo, error)                                       `perm:"admin"`
		WalletImport func(ctx context.Context, ki *core.KeyInfo) (core.Address, error)                                         `perm:"admin"`
		WalletDelete func(ctx context.Context, addr core.Address) error                                                        `perm:"admin"`
	}
}

func (c *WalletAPIAdapter) WalletNew(ctx context.Context, keyType core.KeyType) (core.Address, error) {
	return c.Internal.WalletNew(ctx, keyType)
}

func (c *WalletAPIAdapter) WalletHas(ctx context.Context, addr core.Address) (bool, error) {
	return c.Internal.WalletHas(ctx, addr)
}

func (c *WalletAPIAdapter) WalletList(ctx context.Context) ([]core.Address, error) {
	return c.Internal.WalletList(ctx)
}

func (c *WalletAPIAdapter) WalletSign(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error) {
	return c.Internal.WalletSign(ctx, signer, toSign, meta)
}

func (c *WalletAPIAdapter) WalletExport(ctx context.Context, a core.Address) (*core.KeyInfo, error) {
	return c.Internal.WalletExport(ctx, a)
}

func (c *WalletAPIAdapter) WalletImport(ctx context.Context, ki *core.KeyInfo) (core.Address, error) {
	return c.Internal.WalletImport(ctx, ki)
}

func (c *WalletAPIAdapter) WalletDelete(ctx context.Context, addr core.Address) error {
	return c.Internal.WalletDelete(ctx, addr)
}
