package cli

import (
	"errors"
	"fmt"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/msgrouter"
	"github.com/urfave/cli/v2"
)

var (
	ErrParameterMismatch = errors.New("parameter mismatch")
)
var strategyCmd = &cli.Command{
	Name:  "strategy",
	Usage: "Manage logging",
	Subcommands: []*cli.Command{
		strategyTypeList,
		strategyMethodList,
	},
}

var strategyTypeList = &cli.Command{
	Name:  "types",
	Usage: "show all msgTypes",
	Action: func(cctx *cli.Context) error {
		fmt.Println("code\ttype")
		for _, v := range core.MsgEnumPool {
			fmt.Printf("%d\t%s\n", v.Code, v.Name)
		}
		return nil
	},
}
var strategyMethodList = &cli.Command{
	Name:  "methods",
	Usage: "show all methods (index are used for counting only)",
	Action: func(cctx *cli.Context) error {
		fmt.Println("index\tmethod")
		for k, v := range msgrouter.MethodNameList {
			fmt.Printf("%d\t%s\n", k+1, v)
		}
		return nil
	},
}

var strategyPutMsgTypeTemplate = &cli.Command{
	Name:      "newMsgTypeTemplate",
	Aliases:   []string{"newMTT"},
	Usage:     "create a msgType common template",
	ArgsUsage: "[name, code1 code2 ...]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return ShowHelp(cctx, ErrParameterMismatch)
		}
		/*		name := cctx.Args().First()
				codesMap := make(map[core.MsgEnum]struct{})
				for _, arg := range cctx.Args().Slice()[1:] {
					number, err := strconv.Atoi(arg)
					if err != nil {
						return errors.New("code must be the number")
					}
					code, err := core.MsgEnumFromInt(number)
					if err != nil {
						return err
					}
					codesMap[code] = struct{}{}
				}
				var sumCode core.MsgEnum
				for k, _ := range codesMap {
					sumCode += k
				}

		*/
		return nil
	},
}
