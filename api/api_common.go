package api

import (
	"context"
	"github.com/ipfs-force-community/venus-wallet/api/permission"
	"github.com/ipfs-force-community/venus-wallet/common"
)

var _ common.ICommon = &CommonAuth{}

// common API permissions constraints
type CommonAuth struct {
	Internal struct {
		AuthVerify  func(ctx context.Context, token string) ([]permission.Permission, error) `perm:"read"`
		AuthNew     func(ctx context.Context, perms []permission.Permission) ([]byte, error) `perm:"admin"`
		Version     func(context.Context) (common.Version, error)                            `perm:"read"`
		LogList     func(context.Context) ([]string, error)                                  `perm:"write"`
		LogSetLevel func(context.Context, string, string) error                              `perm:"write"`
	}
}

func (c *CommonAuth) AuthVerify(ctx context.Context, token string) ([]permission.Permission, error) {
	return c.Internal.AuthVerify(ctx, token)
}

func (c *CommonAuth) AuthNew(ctx context.Context, perms []permission.Permission) ([]byte, error) {
	return c.Internal.AuthNew(ctx, perms)
}

// Version implements API.Version
func (c *CommonAuth) Version(ctx context.Context) (common.Version, error) {
	return c.Internal.Version(ctx)
}

func (c *CommonAuth) LogList(ctx context.Context) ([]string, error) {
	return c.Internal.LogList(ctx)
}

func (c *CommonAuth) LogSetLevel(ctx context.Context, group, level string) error {
	return c.Internal.LogSetLevel(ctx, group, level)
}
