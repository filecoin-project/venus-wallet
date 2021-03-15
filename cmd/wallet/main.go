package main

import (
	"context"
	localCli "github.com/ipfs-force-community/venus-wallet/cli"
	"github.com/ipfs-force-community/venus-wallet/cli/helper"
	main2 "github.com/ipfs-force-community/venus-wallet/cmd"
	loclog "github.com/ipfs-force-community/venus-wallet/log"
	"github.com/ipfs-force-community/venus-wallet/middleware"
	"github.com/ipfs-force-community/venus-wallet/version"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli/v2"
	"go.opencensus.io/trace"
	"os"
)

func main() {
	loclog.SetupLogLevels()
	local := []*cli.Command{
		main2.RunCmd,
	}
	jaeger := middleware.SetupJaegerTracing("venus-wallet")
	defer func() {
		if jaeger != nil {
			jaeger.Flush()
		}
	}()
	for _, cmd := range local {
		cmd := cmd
		originBefore := cmd.Before
		cmd.Before = func(cctx *cli.Context) error {
			trace.UnregisterExporter(jaeger)
			jaeger = middleware.SetupJaegerTracing("venus/" + cmd.Name)

			if originBefore != nil {
				return originBefore(cctx)
			}
			return nil
		}
	}
	ctx, span := trace.StartSpan(context.Background(), "/cli")
	defer span.End()

	app := &cli.App{
		Name:    "venus remote-wallet",
		Usage:   "",
		Version: version.UserVersion,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "repo",
				EnvVars: []string{"VENUS_WALLET_PATH"},
				Hidden:  true,
				Value:   "~/.venus_wallet",
			},
		},

		Commands: append(local, localCli.Commands...),
	}
	app.Setup()
	app.Metadata["traceContext"] = ctx

	if err := app.Run(os.Args); err != nil {
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeFailedPrecondition,
			Message: err.Error(),
		})
		_, ok := err.(*helper.ErrCmdFailed)
		if ok {
			log.Debugf("%+v", err)
		} else {
			log.Warnf("%+v", err)
		}
		os.Exit(1)
	}
}
