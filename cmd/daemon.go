package cmd

import (
	"context"

	address "github.com/filecoin-project/go-address"
	homedir "github.com/mitchellh/go-homedir"
	cli "github.com/urfave/cli/v2"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/venus-wallet/api"
	"github.com/filecoin-project/venus-wallet/build"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/filemgr"
	"github.com/filecoin-project/venus-wallet/middleware"
	"github.com/filecoin-project/venus-wallet/version"
)

type cmd = string

const (
	cmdNetwork cmd = "network"
	//cmdAPI     cmd = "api"
	cmdRepo cmd = "repo"
	//cmdKeyStore cmd = "keystore"
	cmdPwd             cmd = "password"
	cmdGatewayAPI      cmd = "gateway-api"
	cmdGatewayToken    cmd = "gateway-token"
	cmdSupportAccounts cmd = "support-accounts"
)

// DaemonCmd is the `go-lotus daemon` command
var RunCmd = &cli.Command{
	Name:  "run",
	Usage: "Start a venus wallet process",
	Flags: []cli.Flag{
		//	&cli.StringFlag{Name: cmdAPI, Value: "5678"},
		&cli.StringFlag{Name: cmdNetwork, Value: ""},
		&cli.StringFlag{Name: cmdPwd, Value: "", Aliases: []string{"pwd"}},
		&cli.StringSliceFlag{Name: cmdGatewayAPI},
		&cli.StringFlag{Name: cmdGatewayToken, Value: ""},
		&cli.StringSliceFlag{Name: cmdSupportAccounts},
	},
	Action: func(cctx *cli.Context) error {
		ctx, _ := tag.New(context.Background(), tag.Insert(middleware.Version, version.BuildVersion))
		dir, err := homedir.Expand(cctx.String(cmdRepo))
		if err != nil {
			log.Warnw("could not expand repo location", "error", err)
		} else {
			log.Infof("wallet repo: %s", dir)
		}
		op := &filemgr.OverrideParams{
			GatewayAPI:      cctx.StringSlice(cmdGatewayAPI),
			GatewayToken:    cctx.String(cmdGatewayToken),
			SupportAccounts: cctx.StringSlice(cmdSupportAccounts),
		}
		r, err := filemgr.NewFS(cctx.String(cmdRepo), op)
		if err != nil {
			return xerrors.Errorf("opening fs repo: %w", err)
		}
		core.WalletStrategyLevel = r.Config().Strategy.Level
		secret, err := r.APISecret()
		if err != nil {
			return xerrors.Errorf("read secret failed: %w", err)
		}
		var fullAPI api.IFullAPI

		if cctx.String(cmdNetwork) == "test" {
			address.CurrentNetwork = address.Testnet
		} else {
			address.CurrentNetwork = address.Mainnet
		}
		stop, err := build.New(ctx,
			build.FullAPIOpt(&fullAPI),
			build.WalletOpt(r, cctx.String(cmdPwd)),
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
