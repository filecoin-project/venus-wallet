package api

import (
	"context"
	"github.com/filecoin-project/venus-wallet/wallet_event"
)

var _ wallet_event.IWalletEventAPI = &WalletEventAPIAdapter{}

type WalletEventAPIAdapter struct {
	Internal struct {
		AddSupportAccount func(ctx context.Context, supportAccount string) error `perm:"admin"`
	}
}

func (w WalletEventAPIAdapter) AddSupportAccount(ctx context.Context, supportAccount string) error {
	return w.Internal.AddSupportAccount(ctx, supportAccount)
}
