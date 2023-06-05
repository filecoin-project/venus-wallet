package cli

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/filecoin-project/venus-auth/core"

	"github.com/filecoin-project/venus-wallet/cli/helper"
)

var authCmd = &cli.Command{
	Name:  "auth",
	Usage: "Manage RPC permissions",
	Subcommands: []*cli.Command{
		authApiInfoToken,
	},
}

var authApiInfoToken = &cli.Command{
	Name:  "api-info",
	Usage: "Get token with API info required to connect to this node",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "perm",
			Usage: "permission to assign to the token, one of: read, write, sign, admin",
		},
	},

	Action: func(cctx *cli.Context) error {
		api, closer, err := helper.GetAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		ctx := helper.ReqContext(cctx)

		if !cctx.IsSet("perm") {
			return errors.New("--perm flag not set")
		}

		allPermissions := core.AdaptOldStrategy(core.PermAdmin)
		perm := cctx.String("perm")
		idx := 0
		for i, p := range allPermissions {
			if perm == p {
				idx = i + 1
			}
		}

		if idx == 0 {
			return fmt.Errorf("--perm flag has to be one of: %s", allPermissions)
		}

		// slice on [:idx] so for example: 'sign' gives you [read, write, sign]
		token, err := api.AuthNew(ctx, allPermissions[:idx])
		if err != nil {
			return err
		}

		apiInfo, err := helper.GetAPIInfo(cctx)
		if err != nil {
			return fmt.Errorf("could not get API info: %w", err)
		}
		fmt.Printf("%s:%s\n", string(token), apiInfo.Addr)
		return nil
	},
}
