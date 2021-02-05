package main

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/venus-wallet/build"
	"github.com/ipfs-force-community/venus-wallet/node/config"
	"github.com/ipfs-force-community/venus-wallet/node/modules/dtypes"
	"github.com/mitchellh/go-homedir"
	"github.com/multiformats/go-multiaddr"
	"github.com/urfave/cli/v2"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"golang.org/x/xerrors"
	"path"

	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/metrics"
	"github.com/ipfs-force-community/venus-wallet/node"
	"github.com/ipfs-force-community/venus-wallet/node/repo"
)

// DaemonCmd is the `go-lotus daemon` command
var DaemonCmd = &cli.Command{
	Name:  "daemon",
	Usage: "Start a lotus daemon process",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "api", Value: "5678"},
		&cli.StringFlag{Name: "keystore", Aliases: []string{"ks"}},
		&cli.StringFlag{Name: "network", Value: ""},
	},
	Action: func(cctx *cli.Context) error {
		ctx, _ := tag.New(context.Background(), tag.Insert(metrics.Version, build.BuildVersion))
		dir, err := homedir.Expand(cctx.String("repo"))
		if err != nil {
			log.Warnw("could not expand repo location", "error", err)
		} else {
			log.Infof("lotus repo: %s", dir)
		}
		r, err := repo.NewFS(cctx.String("repo"))
		if err != nil {
			return xerrors.Errorf("opening fs repo: %w", err)
		}

		if err := r.Init(repo.FullNode); err != nil && err != repo.ErrRepoExists {
			return xerrors.Errorf("repo init error: %w", err)
		}

		if !cctx.IsSet("keystore") {
			_ = cctx.Set("keystore", path.Join(dir, "keystore.sqlit"))
		}

		var api api.FullNode

		stop, err := node.New(ctx,
			node.Override(node.SetNet, func() {
				address.CurrentNetwork = address.Mainnet
			}),
			node.FullAPI(&api),
			node.ApplyIf(func(s *node.Settings) bool { return !cctx.Bool("bootstrap") },
				node.Local(),
			),
			node.Online(),
			node.Repo(r),
			node.ApplyIf(func(s *node.Settings) bool { return cctx.IsSet("api") },
				node.Override(node.SetApiEndpointKey, func(lr repo.LockedRepo) error {
					apima, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/" +
						cctx.String("api"))
					if err != nil {
						return err
					}
					return lr.SetAPIEndpoint(apima)
				})),
			node.Override(new(*config.DbCfg),
				&config.DbCfg{Conn: cctx.String("keystore"), Type: "sqlite", DebugMode: true}),
			node.ApplyIf(func(s *node.Settings) bool { return cctx.IsSet("network") },
				node.Override(new(dtypes.NetworkName), dtypes.NetworkName(cctx.String("network")))),
		)
		if err != nil {
			return xerrors.Errorf("initializing node: %w", err)
		}

		// Register all metric views
		if err = view.Register(
			metrics.DefaultViews...,
		); err != nil {
			log.Fatalf("Cannot register the view: %v", err)
		}

		// Set the metric to one so it is published to the exporter
		stats.Record(ctx, metrics.LotusInfo.M(1))

		endpoint, err := r.APIEndpoint()
		if err != nil {
			return xerrors.Errorf("getting api endpoint: %w", err)
		}

		// TODO: properly parse api endpoint (or make it a URL)
		return ServeRPC(api, stop, endpoint)
	},
}
