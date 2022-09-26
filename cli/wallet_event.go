package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/filecoin-project/venus-wallet/cli/helper"
)

var supportCmds = &cli.Command{
	Name:      "support",
	Usage:     "Add an account that can be signed with the private key",
	ArgsUsage: "<account>",
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
