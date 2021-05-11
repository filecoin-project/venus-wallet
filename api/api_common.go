package api

import (
	"context"
	"github.com/filecoin-project/venus-wallet/api/permission"
	"github.com/filecoin-project/venus-wallet/common"
)

var _ common.ICommon = &CommonAPIAdapter{}

// common API permissions constraints
type CommonAPIAdapter struct {
	Internal struct {
		AuthVerify  func(ctx context.Context, token string) ([]permission.Permission, error) `perm:"read"`
		AuthNew     func(ctx context.Context, perms []permission.Permission) ([]byte, error) `perm:"admin"`
		Version     func(context.Context) (common.Version, error)                            `perm:"read"`
		LogList     func(context.Context) ([]string, error)                                  `perm:"write"`
		LogSetLevel func(context.Context, string, string) error                              `perm:"write"`
	}
}

func (c *CommonAPIAdapter) AuthVerify(ctx context.Context, token string) ([]permission.Permission, error) {
	return c.Internal.AuthVerify(ctx, token)
}

func (c *CommonAPIAdapter) AuthNew(ctx context.Context, perms []permission.Permission) ([]byte, error) {
	return c.Internal.AuthNew(ctx, perms)
}

// Version implements API.Version
func (c *CommonAPIAdapter) Version(ctx context.Context) (common.Version, error) {
	return c.Internal.Version(ctx)
}

func (c *CommonAPIAdapter) LogList(ctx context.Context) ([]string, error) {
	return c.Internal.LogList(ctx)
}

func (c *CommonAPIAdapter) LogSetLevel(ctx context.Context, group, level string) error {
	return c.Internal.LogSetLevel(ctx, group, level)
}
