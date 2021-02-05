package modules

import (
	"context"

	"go.uber.org/fx"

	"github.com/ipfs-force-community/venus-wallet/node/modules/dtypes"
	"github.com/ipfs-force-community/venus-wallet/node/repo"
	"github.com/filecoin-project/lotus/chain/types"
)

func LockedRepo(lr repo.LockedRepo) func(lc fx.Lifecycle) repo.LockedRepo {
	return func(lc fx.Lifecycle) repo.LockedRepo {
		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return lr.Close()
			},
		})

		return lr
	}
}

func KeyStore(lr repo.LockedRepo) (types.KeyStore, error) {
	return lr.KeyStore()
}

func Datastore(r repo.LockedRepo) (dtypes.MetadataDS, error) {
	return r.Datastore("/metadata")
}
