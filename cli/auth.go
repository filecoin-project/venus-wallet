package cli

import (
	"fmt"

	"github.com/filecoin-project/venus-wallet/cli/helper"
	"github.com/filecoin-project/venus/venus-shared/api/permission"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

var authCmd = &cli.Command{
	Name:  "auth",
	Usage: "Manage RPC permissions",
	Subcommands: []*cli.Command{
		authCreateAdminToken,
		authApiInfoToken,
	},
}

var authCreateAdminToken = &cli.Command{
	Name:  "create-token",
	Usage: "Create token",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "perm",
			Usage: "permission to assign to the token, one of: read, write, sign, admin",
		},
	},

	Action: func(cctx *cli.Context) error {
		napi, closer, err := helper.GetAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		ctx := helper.ReqContext(cctx)

		if !cctx.IsSet("perm") {
			return xerrors.New("--perm flag not set")
		}

		perm := cctx.String("perm")
		idx := 0
		for i, p := range permission.AllPermissions {
			if perm == p {
				idx = i + 1
			}
		}

		if idx == 0 {
			return fmt.Errorf("--perm flag has to be one of: %s", permission.AllPermissions)
		}

		// slice on [:idx] so for example: 'sign' gives you [read, write, sign]
		token, err := napi.AuthNew(ctx, permission.AllPermissions[:idx])
		if err != nil {
			return err
		}

		// TODO: Log in audit log when it is implemented

		fmt.Println(string(token))
		return nil
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
		napi, closer, err := helper.GetAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		ctx := helper.ReqContext(cctx)

		if !cctx.IsSet("perm") {
			return xerrors.New("--perm flag not set")
		}

		perm := cctx.String("perm")
		idx := 0
		for i, p := range permission.AllPermissions {
			if perm == p {
				idx = i + 1
			}
		}

		if idx == 0 {
			return fmt.Errorf("--perm flag has to be one of: %s", permission.AllPermissions)
		}

		// slice on [:idx] so for example: 'sign' gives you [read, write, sign]
		token, err := napi.AuthNew(ctx, permission.AllPermissions[:idx])
		if err != nil {
			return err
		}

		ainfo, err := helper.GetAPIInfo(cctx)
		if err != nil {
			return xerrors.Errorf("could not get API info: %w", err)
		}
		fmt.Printf("%s:%s\n", string(token), ainfo.Addr)
		return nil
	},
}
