package build

import (
	"context"
	"fmt"

	"github.com/asaskevich/EventBus"
	"github.com/filecoin-project/venus-wallet/api"
	"github.com/filecoin-project/venus-wallet/common"
	"github.com/filecoin-project/venus-wallet/config"
	"github.com/filecoin-project/venus-wallet/filemgr"
	"github.com/filecoin-project/venus-wallet/storage"
	"github.com/filecoin-project/venus-wallet/storage/sqlite"
	"github.com/filecoin-project/venus-wallet/storage/wallet"
	"github.com/filecoin-project/venus-wallet/wallet_event"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/ipfs-force-community/sophon-gateway/types"
	"go.uber.org/fx"
	"gorm.io/gorm"

	wallet_api "github.com/filecoin-project/venus/venus-shared/api/wallet"
)

// special is a type used to give keys to modules which
//
//	can't really be identified by the returned type
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

func WalletOpt(repo filemgr.Repo, walletPwd string) Option {
	c := repo.Config()
	return Options(
		Override(new(filemgr.Repo), repo),
		Override(new(*config.DBConfig), c.DB),
		Override(new(EventBus.Bus), EventBus.New),
		Override(new(*gorm.DB), sqlite.NewDB),
		Override(new(*config.CryptoFactor), c.Factor),
		Override(new(storage.KeyMiddleware), storage.NewKeyMiddleware),
		Override(new(storage.KeyStore), sqlite.NewKeyStore),
		Override(new(*config.SignRecorderConfig), c.SignRecorder),
		Override(new(storage.IRecorder), sqlite.NewSqliteRecorder),
		Override(new(wallet.GetPwdFunc), func() wallet.GetPwdFunc {
			return func() string {
				return walletPwd
			}
		}),
		Override(new(wallet.ISignMsgFilter), func() wallet.ISignMsgFilter {
			return wallet.NewSignFilter(c.SignFilter)
		}),
		Override(new(wallet_api.ILocalWallet), wallet.NewWallet),

		Override(new(types.IWalletHandler), From(new(wallet_api.ILocalWallet))),
		Override(new(*config.APIRegisterHubConfig), c.APIRegisterHub),
		Override(new(wallet_event.IAPIRegisterHub), wallet_event.NewAPIRegisterHub),
		Override(new(wallet_api.IWalletEvent), wallet_event.NewWalletEventAPI),
	)
}

func CommonOpt(alg *jwt.HMACSHA) Option {
	return Options(
		Override(new(*jwt.HMACSHA), alg),
		Override(new(common.ICommon), From(new(common.Common))),
	)
}

func FullAPIOpt(out *api.IFullAPI) Option {
	return func(s *Settings) error {
		resAPI := &api.FullAPI{}
		s.invokes[ExtractApiKey] = fx.Populate(resAPI)
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
		return nil, fmt.Errorf("applying node options failed: %w", err)
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
		return nil, fmt.Errorf("starting node: %w", err)
	}
	return app.Stop, nil
}
