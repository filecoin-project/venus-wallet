package apistruct

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	api2 "github.com/filecoin-project/lotus/api"

	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/filecoin-project/lotus/chain/types"
)

// All permissions are listed in permissioned.go
var _ = AllPermissions

type CommonStruct struct {
	Internal struct {
		AuthVerify  func(ctx context.Context, token string) ([]api.Permission, error) `perm:"read"`
		AuthNew     func(ctx context.Context, perms []api.Permission) ([]byte, error) `perm:"admin"`
		Version     func(context.Context) (api.Version, error)                        `perm:"read"`
		LogList     func(context.Context) ([]string, error)                           `perm:"write"`
		LogSetLevel func(context.Context, string, string) error                       `perm:"write"`
	}
}

// FullNodeStruct implements API passing calls to user-provided function values.
type FullNodeStruct struct {
	CommonStruct
	Internal struct {
		// remote wallet api
		WalletNew    func(context.Context, types.KeyType) (address.Address, error)                           `perm:"admin"`
		WalletHas    func(context.Context, address.Address) (bool, error)                                    `perm:"write"`
		WalletList   func(context.Context) ([]address.Address, error)                                        `perm:"write"`
		WalletSign   func(context.Context, address.Address, []byte, api2.MsgMeta) (*crypto.Signature, error) `perm:"sign"`
		WalletExport func(context.Context, address.Address) (*types.KeyInfo, error)                          `perm:"admin"`
		WalletImport func(context.Context, *types.KeyInfo) (address.Address, error)                          `perm:"admin"`
		WalletDelete func(context.Context, address.Address) error                                            `perm:"admin"`
	}
}

func (c *CommonStruct) AuthVerify(ctx context.Context, token string) ([]api.Permission, error) {
	return c.Internal.AuthVerify(ctx, token)
}

func (c *CommonStruct) AuthNew(ctx context.Context, perms []api.Permission) ([]byte, error) {
	return c.Internal.AuthNew(ctx, perms)
}

// Version implements API.Version
func (c *CommonStruct) Version(ctx context.Context) (api.Version, error) {
	return c.Internal.Version(ctx)
}

func (c *CommonStruct) LogList(ctx context.Context) ([]string, error) {
	return c.Internal.LogList(ctx)
}

func (c *CommonStruct) LogSetLevel(ctx context.Context, group, level string) error {
	return c.Internal.LogSetLevel(ctx, group, level)
}

func (c *FullNodeStruct) WalletNew(ctx context.Context, keyType types.KeyType) (address.Address, error) {
	return c.Internal.WalletNew(ctx, keyType)
}

func (c *FullNodeStruct) WalletHas(ctx context.Context, addr address.Address) (bool, error) {
	return c.Internal.WalletHas(ctx, addr)
}

func (c *FullNodeStruct) WalletList(ctx context.Context) ([]address.Address, error) {
	return c.Internal.WalletList(ctx)
}

func (c *FullNodeStruct) WalletSign(ctx context.Context, signer address.Address, toSign []byte, meta api2.MsgMeta) (*crypto.Signature, error) {
	return c.Internal.WalletSign(ctx, signer, toSign, meta)
}

func (c *FullNodeStruct) WalletExport(ctx context.Context, a address.Address) (*types.KeyInfo, error) {
	return c.Internal.WalletExport(ctx, a)
}

func (c *FullNodeStruct) WalletImport(ctx context.Context, ki *types.KeyInfo) (address.Address, error) {
	return c.Internal.WalletImport(ctx, ki)
}

func (c *FullNodeStruct) WalletDelete(ctx context.Context, addr address.Address) error {
	return c.Internal.WalletDelete(ctx, addr)
}

var _ api.Common = &CommonStruct{}
var _ api.FullNode = &FullNodeStruct{}
