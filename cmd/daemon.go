package main

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/build"
	"github.com/ipfs-force-community/venus-wallet/filemgr"
	"github.com/ipfs-force-community/venus-wallet/middleware"
	"github.com/ipfs-force-community/venus-wallet/version"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"golang.org/x/xerrors"
)

type cmd = string

const (
	cmdNetwork cmd = "network"
	cmdAPI     cmd = "api"
	cmdRepo    cmd = "repo"
	//cmdKeyStore cmd = "keystore"
)

// DaemonCmd is the `go-lotus daemon` command
var RunCmd = &cli.Command{
	Name:  "run",
	Usage: "Start a venus wallet process",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: cmdAPI, Value: "5678"},
		&cli.StringFlag{Name: cmdNetwork, Value: ""},
	},
	Action: func(cctx *cli.Context) error {
		ctx, _ := tag.New(context.Background(), tag.Insert(middleware.Version, version.BuildVersion))
		dir, err := homedir.Expand(cctx.String(cmdRepo))
		if err != nil {
			log.Warnw("could not expand repo location", "error", err)
		} else {
			log.Infof("wallet repo: %s", dir)
		}
		apiListen := ""
		if cctx.IsSet("api") {
			apiListen = "/ip4/0.0.0.0/tcp/" + cctx.String("api")
		}
		op := &filemgr.OverrideParams{
			API: apiListen,
		}
		r, err := filemgr.NewFS(cctx.String(cmdRepo), op)
		if err != nil {
			return xerrors.Errorf("opening fs repo: %w", err)
		}

		secret, err := r.APISecret()
		if err != nil {
			return xerrors.Errorf("read secret failed: %w", err)
		}
		var fullAPI api.IFullAPI

		stop, err := build.New(ctx,
			build.Override(build.SetNet, func() {
				address.CurrentNetwork = address.Mainnet
			}),
			build.FullAPIOpt(&fullAPI),
			build.WalletOpt(r.Config()),
			build.CommonOpt(secret),
			build.ApplyIf(func(s *build.Settings) bool { return cctx.IsSet(cmdNetwork) },
				build.Override(new(build.NetworkName), build.NetworkName(cctx.String(cmdNetwork)))),
		)
		if err != nil {
			return xerrors.Errorf("initializing node: %w", err)
		}

		// Register all metric views
		if err = view.Register(
			middleware.DefaultViews...,
		); err != nil {
			log.Fatalf("Cannot register the view: %v", err)
		}

		// Set the metric to one so it is published to the exporter
		stats.Record(ctx, middleware.VenusInfo.M(1))

		endpoint, err := r.APIEndpoint()
		if err != nil {
			return xerrors.Errorf("getting api endpoint: %w", err)
		}

		// TODO: properly parse api endpoint (or make it a URL)
		return ServeRPC(fullAPI, stop, endpoint)
	},
}
