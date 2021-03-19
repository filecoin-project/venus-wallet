package api

import (
	"context"
	"github.com/ipfs-force-community/venus-wallet/storage"
)

var _ storage.IWalletLock = &WalletLockAuth{}

type WalletLockAuth struct {
	Internal struct {
		SetPassword func(ctx context.Context, password string) error `perm:"admin" local:"required"`
		Unlock      func(ctx context.Context, password string) error `perm:"admin" local:"required"`
		Lock        func(ctx context.Context, password string) error `perm:"admin" local:"required"`
		LockState   func(ctx context.Context) bool                   `perm:"admin"`
	}
}

func (c *WalletLockAuth) SetPassword(ctx context.Context, password string) error {
	return c.Internal.SetPassword(ctx, password)
}
func (c *WalletLockAuth) Unlock(ctx context.Context, password string) error {
	return c.Internal.Unlock(ctx, password)
}
func (c *WalletLockAuth) Lock(ctx context.Context, password string) error {
	return c.Internal.Lock(ctx, password)
}

func (c *WalletLockAuth) LockState(ctx context.Context) bool {
	return c.Internal.LockState(ctx)
}
