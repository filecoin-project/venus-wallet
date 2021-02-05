package node

import (
	"context"
	"errors"
	"github.com/ipfs-force-community/venus-wallet/api"
	_ "github.com/ipfs-force-community/venus-wallet/lib/sigs/bls"
	_ "github.com/ipfs-force-community/venus-wallet/lib/sigs/secp" // 为了调用init函数
	"github.com/ipfs-force-community/venus-wallet/node/config"
	"github.com/ipfs-force-community/venus-wallet/node/impl"
	"github.com/ipfs-force-community/venus-wallet/node/impl/force/db_proc"
	"github.com/ipfs-force-community/venus-wallet/node/modules"
	"github.com/ipfs-force-community/venus-wallet/node/modules/dtypes"
	"github.com/ipfs-force-community/venus-wallet/node/modules/helpers"
	"github.com/ipfs-force-community/venus-wallet/node/repo"
	"github.com/filecoin-project/lotus/chain/types"
	logging "github.com/ipfs/go-log"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/fx"
	"golang.org/x/xerrors"
)

var log = logging.Logger("builder")

// special is a type used to give keys to modules which
//  can't really be identified by the returned type
type special struct{ id int }

type invoke int

// nolint:golint
const (
	PstoreAddSelfKeysKey = invoke(iota)
	ExtractApiKey
	SetNet
	SetApiEndpointKey
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

	nodeType repo.RepoType

	Online bool // Online option applied
	Config bool // Config option applied

	Local bool
}

func defaults() []Option {
	return []Option{
		Override(new(helpers.MetricsCtx), context.Background),
	}
}

func isType(t repo.RepoType) func(s *Settings) bool {
	return func(s *Settings) bool { return s.nodeType == t }
}

func Local() Option {
	return Options(
		// make sure that local is applied before Config.
		// This is important ***
		func(s *Settings) error { s.Local = true; return nil },
	)
}

func Online() Option {
	return Options(
		func(s *Settings) error { s.Online = true; return nil },
		ApplyIf(func(s *Settings) bool { return s.Config },
			Error(errors.New("the Online option must be set before Config option")),
		),
		// Full node
		ApplyIf(isType(repo.FullNode),
			Override(new(db_proc.DbProcInterface), newDbProc),
		),
	)
}

// Config sets up constructors based on the provided Config
func ConfigCommon(cfg *config.Common) Option {
	return Options(
		func(s *Settings) error { s.Config = true; return nil },
		Override(new(dtypes.APIEndpoint), func() (dtypes.APIEndpoint, error) {
			return multiaddr.NewMultiaddr(cfg.API.ListenAddress)
		}),
		Override(SetApiEndpointKey, func(lr repo.LockedRepo, e dtypes.APIEndpoint) error {
			return lr.SetAPIEndpoint(e)
		}),
		ApplyIf(func(s *Settings) bool { return s.Online },
			Override(new(*config.DbCfg), &cfg.DbCfg),
		),
	)
}

func ConfigFullNode(c interface{}) Option {
	cfg, ok := c.(*config.FullNode)
	if !ok {
		return Error(xerrors.Errorf("invalid config from repo, got: %T", c))
	}

	return Options(
		ConfigCommon(&cfg.Common),
	)
}

func Repo(r repo.Repo) Option {
	return func(settings *Settings) error {
		lr, err := r.Lock(settings.nodeType)
		if err != nil {
			return err
		}
		c, err := lr.Config()
		if err != nil {
			return err
		}

		return Options(
			Override(new(repo.LockedRepo), modules.LockedRepo(lr)), // module handles closing
			Override(new(dtypes.MetadataDS), modules.Datastore),
			Override(new(types.KeyStore), modules.KeyStore),
			Override(new(*dtypes.APIAlg), modules.APISecret),
			ApplyIf(isType(repo.FullNode), ConfigFullNode(c)),
		)(settings)
	}
}

func FullAPI(out *api.FullNode) Option {
	return func(s *Settings) error {
		resAPI := &impl.FullNodeAPI{}
		s.invokes[ExtractApiKey] = fx.Extract(resAPI)
		*out = resAPI
		return nil
	}
}

type StopFunc func(context.Context) error

// New builds and starts new Filecoin node
func New(ctx context.Context, opts ...Option) (StopFunc, error) {
	settings := Settings{
		modules:  map[interface{}]fx.Option{},
		invokes:  make([]fx.Option, _nInvokes),
		nodeType: repo.FullNode,
	}

	// apply module options in the right order
	if err := Options(Options(defaults()...), Options(opts...))(&settings); err != nil {
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

	// TODO: we probably should have a 'firewall' for Closing signal
	//  on this context, and implement closing logic through lifecycles
	//  correctly
	if err := app.Start(ctx); err != nil {
		// comment fx.NopLogger few lines above for easier debugging
		return nil, xerrors.Errorf("starting node: %w", err)
	}

	return app.Stop, nil
}

func newDbProc(cfg *config.DbCfg) (db_proc.DbProcInterface, error) {
	if cfg != nil && cfg.Conn != "" {
		return db_proc.NewDbProc(cfg)
	}
	return nil, xerrors.New("mysql conn do not set")
}
