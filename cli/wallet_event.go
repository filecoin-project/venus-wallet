package cli

import (
	"fmt"

	"github.com/filecoin-project/venus-wallet/cli/helper"
	"github.com/urfave/cli/v2"
)

var supportCmds = &cli.Command{
	Name:      "support",
	Aliases:   []string{"support"},
	Usage:     "tell upstream which account to support",
	ArgsUsage: "account",
	Action: func(cctx *cli.Context) error {
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		if cctx.NArg() != 1 {
			return fmt.Errorf("must specify account to support")
		}
		return api.AddSupportAccount(ctx, cctx.Args().Get(0))
	},
}
