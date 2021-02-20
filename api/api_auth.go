package api

import (
	"context"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/storage"
)

var _ ICommon = &CommonAuth{}
var _ IFullAPI = &ServerAuth{}
var _ storage.IWallet = &WalletAuth{}

type ServerAuth struct {
	CommonAuth
	WalletAuth
}

type CommonAuth struct {
	internal struct {
		AuthVerify  func(ctx context.Context, token string) ([]Permission, error) `perm:"read"`
		AuthNew     func(ctx context.Context, perms []Permission) ([]byte, error) `perm:"admin"`
		Version     func(context.Context) (Version, error)                        `perm:"read"`
		LogList     func(context.Context) ([]string, error)                       `perm:"write"`
		LogSetLevel func(context.Context, string, string) error                   `perm:"write"`
	}
}

type WalletAuth struct {
	internal struct {
		WalletNew    func(context.Context, core.KeyType) (core.Address, error)                                                 `perm:"admin"`
		WalletHas    func(ctx context.Context, address core.Address) (bool, error)                                             `perm:"write"`
		WalletList   func(ctx context.Context) ([]core.Address, error)                                                         `perm:"write"`
		WalletSign   func(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error) `perm:"sign"`
		WalletExport func(ctx context.Context, addr core.Address) (*core.KeyInfo, error)                                       `perm:"admin"`
		WalletImport func(context.Context, *core.KeyInfo) (core.Address, error)                                                `perm:"admin"`
		WalletDelete func(context.Context, core.Address) error                                                                 `perm:"admin"`
	}
}

func (c *CommonAuth) AuthVerify(ctx context.Context, token string) ([]Permission, error) {
	return c.internal.AuthVerify(ctx, token)
}

func (c *CommonAuth) AuthNew(ctx context.Context, perms []Permission) ([]byte, error) {
	return c.internal.AuthNew(ctx, perms)
}

// Version implements API.Version
func (c *CommonAuth) Version(ctx context.Context) (Version, error) {
	return c.internal.Version(ctx)
}

func (c *CommonAuth) LogList(ctx context.Context) ([]string, error) {
	return c.internal.LogList(ctx)
}

func (c *CommonAuth) LogSetLevel(ctx context.Context, group, level string) error {
	return c.internal.LogSetLevel(ctx, group, level)
}

func (c *WalletAuth) WalletNew(ctx context.Context, keyType core.KeyType) (core.Address, error) {
	return c.internal.WalletNew(ctx, keyType)
}

func (c *WalletAuth) WalletHas(ctx context.Context, addr core.Address) (bool, error) {
	return c.internal.WalletHas(ctx, addr)
}

func (c *WalletAuth) WalletList(ctx context.Context) ([]core.Address, error) {
	return c.internal.WalletList(ctx)
}

func (c *WalletAuth) WalletSign(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error) {
	return c.internal.WalletSign(ctx, signer, toSign, meta)
}

func (c *WalletAuth) WalletExport(ctx context.Context, a core.Address) (*core.KeyInfo, error) {
	return c.internal.WalletExport(ctx, a)
}

func (c *WalletAuth) WalletImport(ctx context.Context, ki *core.KeyInfo) (core.Address, error) {
	return c.internal.WalletImport(ctx, ki)
}

func (c *WalletAuth) WalletDelete(ctx context.Context, addr core.Address) error {
	return c.internal.WalletDelete(ctx, addr)
}
