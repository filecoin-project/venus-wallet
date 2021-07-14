package api

import (
	"context"

	"github.com/filecoin-project/venus-wallet/storage"
)

var _ storage.IWalletLock = &WalletLockAPIAdapter{}

type WalletLockAPIAdapter struct {
	Internal struct {
		SetPassword    func(ctx context.Context, password string) error `perm:"admin" local:"required"`
		Unlock         func(ctx context.Context, password string) error `perm:"admin" local:"required"`
		Lock           func(ctx context.Context, password string) error `perm:"admin" local:"required"`
		VerifyPassword func(ctx context.Context, password string) error `perm:"admin" local:"required"`
		LockState      func(ctx context.Context) bool                   `perm:"admin"`
	}
}

func (c *WalletLockAPIAdapter) VerifyPassword(ctx context.Context, password string) error {
	return c.Internal.VerifyPassword(ctx, password)
}

func (c *WalletLockAPIAdapter) SetPassword(ctx context.Context, password string) error {
	return c.Internal.SetPassword(ctx, password)
}
func (c *WalletLockAPIAdapter) Unlock(ctx context.Context, password string) error {
	return c.Internal.Unlock(ctx, password)
}
func (c *WalletLockAPIAdapter) Lock(ctx context.Context, password string) error {
	return c.Internal.Lock(ctx, password)
}

func (c *WalletLockAPIAdapter) LockState(ctx context.Context) bool {
	return c.Internal.LockState(ctx)
}
