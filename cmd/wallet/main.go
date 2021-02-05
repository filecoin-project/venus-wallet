package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"
	"go.opencensus.io/trace"

	"github.com/ipfs-force-community/venus-wallet/build"
	lcli "github.com/ipfs-force-community/venus-wallet/cli"
	"github.com/ipfs-force-community/venus-wallet/lib/lotuslog"
	"github.com/ipfs-force-community/venus-wallet/lib/tracing"
	"github.com/ipfs-force-community/venus-wallet/node/repo"
)

func main() {
	lotuslog.SetupLogLevels()

	local := []*cli.Command{
		DaemonCmd,
	}

	jaeger := tracing.SetupJaegerTracing("venus-wallet")
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
			jaeger = tracing.SetupJaegerTracing("lotus/" + cmd.Name)

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
		Version: build.UserVersion,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "repo",
				EnvVars: []string{"VENUS_WALLET_PATH"},
				Hidden:  true,
				Value:   "~/.venus_wallet",
			},
		},

		Commands: append(local, lcli.Commands...),
	}
	app.Setup()
	app.Metadata["traceContext"] = ctx
	app.Metadata["repoType"] = repo.FullNode

	if err := app.Run(os.Args); err != nil {
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeFailedPrecondition,
			Message: err.Error(),
		})
		_, ok := err.(*lcli.ErrCmdFailed)
		if ok {
			log.Debugf("%+v", err)
		} else {
			log.Warnf("%+v", err)
		}
		os.Exit(1)
	}
}
