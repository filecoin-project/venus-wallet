package build

import (
	"context"
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/config"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"github.com/ipfs-force-community/venus-wallet/storage/sqlite"
	"github.com/ipfs-force-community/venus-wallet/storage/wallet"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

// special is a type used to give keys to modules which
//  can't really be identified by the returned type
type special struct {
	id int // nolint
}

type invoke int

// nolint:golint
const (
	PstoreAddSelfKeysKey = invoke(iota)
	ExtractApiKey
	SetNet
	_nInvokes // keep this last
)

type Settings struct {
	// modules is a map of constructors for DI
	//
	// In most cases the index will be a reflect. Type of element returned by
	// the constructor, but for some 'constructors' it's hard to specify what's
	// the return type should be (or the constructor returns fx group)
	modules map[interface{}]fx.Option

	// invokes are separate from modules as they can't be referenced by return
	// type, and must be applied in correct order
	invokes []fx.Option
}

func defaults() []Option {
	return []Option{
		Override(new(MetricsCtx), context.Background),
	}
}

func WalletOpt(c *config.Config) Option {
	return Options(
		Override(new(*config.DBConfig), c.DB),
		Override(new(storage.KeyStore), sqlite.NewSQLiteStorage),
		Override(new(api.IWallet), wallet.NewWallet),
	)
}
func CommonOpt(alg *api.APIAlg) Option {
	return Options(
		Override(new(*api.APIAlg), alg),
		Override(new(api.ICommon), From(new(api.Common))),
	)

}
func FullAPIOpt(out *api.IFullAPI) Option {
	return func(s *Settings) error {
		resAPI := &api.FullAPI{}
		s.invokes[ExtractApiKey] = fx.Extract(resAPI)
		*out = resAPI
		return nil
	}
}

type StopFunc func(context.Context) error

// New builds and starts new FileCoin Wallet
func New(ctx context.Context, opts ...Option) (StopFunc, error) {
	settings := Settings{
		modules: map[interface{}]fx.Option{},
		invokes: make([]fx.Option, _nInvokes),
	}
	// apply module options in the right order
	if err := Options(
		Options(defaults()...),
		Options(opts...),
	)(&settings); err != nil {
		return nil, xerrors.Errorf("applying node options failed: %w", err)
	}
	// gather constructors for fx.Options
	ctors := make([]fx.Option, 0, len(settings.modules))
	for _, opt := range settings.modules {
		ctors = append(ctors, opt)
	}
	// fill holes in invokes for use in fx.Options
	for i, opt := range settings.invokes {
		if opt == nil {
			settings.invokes[i] = fx.Options()
		}
	}
	app := fx.New(
		fx.Options(ctors...),
		fx.Options(settings.invokes...),

		fx.NopLogger,
	)

	if err := app.Start(ctx); err != nil {
		// comment fx.NopLogger few lines above for easier debugging
		return nil, xerrors.Errorf("starting node: %w", err)
	}
	return app.Stop, nil
}
