package api

import (
	"context"
	"github.com/ipfs-force-community/venus-wallet/core"
)

var _ ICommon = &CommonAuth{}
var _ IFullAPI = &ServerAuth{}

type ServerAuth struct {
	CommonAuth
	wallet WalletAuth
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
	WalletNew    func(context.Context, core.KeyType) (core.Address, error)                                                 `perm:"admin"`
	WalletHas    func(ctx context.Context, address core.Address) (bool, error)                                             `perm:"write"`
	WalletList   func(ctx context.Context) ([]core.Address, error)                                                         `perm:"write"`
	WalletSign   func(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error) `perm:"sign"`
	WalletExport func(ctx context.Context, addr core.Address) (*core.KeyInfo, error)                                       `perm:"admin"`
	WalletImport func(context.Context, *core.KeyInfo) (core.Address, error)                                                `perm:"admin"`
	WalletDelete func(context.Context, core.Address) error                                                                 `perm:"admin"`
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

func (c *ServerAuth) WalletNew(ctx context.Context, keyType core.KeyType) (core.Address, error) {
	return c.wallet.WalletNew(ctx, keyType)
}

func (c *ServerAuth) WalletHas(ctx context.Context, addr core.Address) (bool, error) {
	return c.wallet.WalletHas(ctx, addr)
}

func (c *ServerAuth) WalletList(ctx context.Context) ([]core.Address, error) {
	return c.wallet.WalletList(ctx)
}

func (c *ServerAuth) WalletSign(ctx context.Context, signer core.Address, toSign []byte, meta core.MsgMeta) (*core.Signature, error) {
	return c.wallet.WalletSign(ctx, signer, toSign, meta)
}

func (c *ServerAuth) WalletExport(ctx context.Context, a core.Address) (*core.KeyInfo, error) {
	return c.wallet.WalletExport(ctx, a)
}

func (c *ServerAuth) WalletImport(ctx context.Context, ki *core.KeyInfo) (core.Address, error) {
	return c.wallet.WalletImport(ctx, ki)
}

func (c *ServerAuth) WalletDelete(ctx context.Context, addr core.Address) error {
	return c.wallet.WalletDelete(ctx, addr)
}
