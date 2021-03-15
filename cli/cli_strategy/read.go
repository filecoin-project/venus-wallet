package cli_strategy

import (
	"errors"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/ipfs-force-community/venus-wallet/cli/helper"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/msgrouter"
	"github.com/urfave/cli/v2"
	"strconv"
	"strings"
)

var (
	ErrParameterMismatch = errors.New("parameter mismatch")
)

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

var strategyGetMsgTypeTemplate = &cli.Command{
	Name:      "msgTypeTemplate",
	Aliases:   []string{"mtt"},
	Usage:     "show msgTypeTemplate by name",
	ArgsUsage: "[name]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		name := cctx.Args().First()
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		tmp, err := api.GetMsgTypeTemplate(name)
		if err != nil {
			return err
		}
		var codes []string
		linq.From(core.FindCode(tmp.MetaTypes)).SelectT(func(i int) string {
			return strconv.FormatInt(int64(i), 10)
		}).ToSlice(&codes)
		fmt.Println(strings.Join(codes, ","))
		return nil
	},
}

var strategyGetMethodTemplateByName = &cli.Command{
	Name:      "methodTemplateByName",
	Aliases:   []string{"mt"},
	Usage:     "show methodTemplate by name",
	ArgsUsage: "[name]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		name := cctx.Args().First()
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		tmp, err := api.GetMethodTemplateByName(name)
		if err != nil {
			return err
		}
		fmt.Println(strings.Join(tmp.Methods, ","))
		return nil
	},
}

var strategyGetKeyBindByName = &cli.Command{
	Name:      "keyBind",
	Aliases:   []string{"kb"},
	Usage:     "show keyBind by name",
	ArgsUsage: "[name]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		name := cctx.Args().First()
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		tmp, err := api.GetKeyBindByName(name)
		if err != nil {
			return err
		}
		fmt.Printf("address: %s", tmp.Address)
		var codes []string
		linq.From(core.FindCode(tmp.MetaTypes)).SelectT(func(i int) string {
			return strconv.FormatInt(int64(i), 10)
		}).ToSlice(&codes)
		fmt.Printf("msgTypes: %s", strings.Join(codes, ","))
		fmt.Printf("methods: %s", strings.Join(tmp.Methods, ","))
		return nil
	},
}

var strategyGetKeyBinds = &cli.Command{
	Name:      "KeyBinds",
	Aliases:   []string{"kbs"},
	Usage:     "show keyBinds by address",
	ArgsUsage: "[address]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		address := cctx.Args().First()
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		arr, err := api.GetKeyBinds(address)
		if err != nil {
			return err
		}
		for k, v := range arr {
			var codes []string
			linq.From(core.FindCode(v.MetaTypes)).SelectT(func(i int) string {
				return strconv.FormatInt(int64(i), 10)
			}).ToSlice(&codes)
			fmt.Printf("num\t:%d\n", k+1)
			fmt.Printf("kbName\t:%s\n", v.Name)
			fmt.Printf("msgTypes\t: %s\n", strings.Join(codes, ","))
			fmt.Printf("methods\t: %s\n\n", strings.Join(v.Methods, ","))
		}
		return nil
	},
}

var strategyGetGroupByName = &cli.Command{
	Name:      "group",
	Aliases:   []string{"kbs"},
	Usage:     "show group by name",
	ArgsUsage: "[name]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		name := cctx.Args().First()
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		group, err := api.GetGroupByName(name)
		if err != nil {
			return err
		}
		for k, v := range group.KeyBinds {
			var codes []string
			linq.From(core.FindCode(v.MetaTypes)).SelectT(func(i int) string {
				return strconv.FormatInt(int64(i), 10)
			}).ToSlice(&codes)
			fmt.Printf("num\t:%d\n", k+1)
			fmt.Printf("kbName\t:%s\n", v.Name)
			fmt.Printf("msgTypes\t: %s\n", strings.Join(codes, ","))
			fmt.Printf("methods\t: %s\n\n", strings.Join(v.Methods, ","))
		}
		return nil
	},
}

var strategyListGroup = &cli.Command{
	Name:      "listGroup",
	Aliases:   []string{"lg"},
	Usage:     "show a range of groups (the element of groups only contain name)",
	ArgsUsage: "[from to]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		from, to, err := helper.ReqFromTo(cctx, 0)
		if err != nil {
			return err
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		groups, err := api.ListGroups(int(from), int(to))
		if err != nil {
			return err
		}
		for k, v := range groups {
			fmt.Printf("%d name\t:%s\n", k+1, v.Name)
		}
		return nil
	},
}

var strategyListKeyBinds = &cli.Command{
	Name:      "listKeyBinds",
	Aliases:   []string{"lkb"},
	Usage:     "show a range of keyBinds (the element of groups only contain name)",
	ArgsUsage: "[from to]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		from, to, err := helper.ReqFromTo(cctx, 0)
		if err != nil {
			return err
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		arr, err := api.ListKeyBinds(int(from), int(to))
		if err != nil {
			return err
		}
		for k, v := range arr {
			var codes []string
			linq.From(core.FindCode(v.MetaTypes)).SelectT(func(i int) string {
				return strconv.FormatInt(int64(i), 10)
			}).ToSlice(&codes)
			fmt.Printf("num\t:%d\n", k+1)
			fmt.Printf("kbName\t:%s\n", v.Name)
			fmt.Printf("msgTypes\t: %s\n", strings.Join(codes, ","))
			fmt.Printf("methods\t: %s\n\n", strings.Join(v.Methods, ","))
		}
		return nil
	},
}

var strategyListMethodTemplates = &cli.Command{
	Name:      "ListMethodTemplates",
	Aliases:   []string{"lmt"},
	Usage:     "show a range of method templates",
	ArgsUsage: "[from to]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		from, to, err := helper.ReqFromTo(cctx, 0)
		if err != nil {
			return err
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		arr, err := api.ListMethodTemplates(int(from), int(to))
		if err != nil {
			return err
		}
		for k, v := range arr {
			fmt.Printf("num\t:%d\n", k+1)
			fmt.Printf("name\t:%s\n", v.Name)
			fmt.Printf("methods\t: %s\n\n", strings.Join(v.Methods, ","))
		}
		return nil
	},
}
var strategyListMsgTypeTemplates = &cli.Command{
	Name:      "ListMsgTypeTemplates",
	Aliases:   []string{"lmtt"},
	Usage:     "show a range of method templates",
	ArgsUsage: "[from to]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		from, to, err := helper.ReqFromTo(cctx, 0)
		if err != nil {
			return err
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		arr, err := api.ListMsgTypeTemplates(int(from), int(to))
		if err != nil {
			return err
		}
		for k, v := range arr {
			var codes []string
			linq.From(core.FindCode(v.MetaTypes)).SelectT(func(i int) string {
				return strconv.FormatInt(int64(i), 10)
			}).ToSlice(&codes)
			fmt.Printf("num\t:%d\n", k+1)
			fmt.Printf("name\t:%s\n", v.Name)
			fmt.Printf("msgTypes\t: %s\n\n", strings.Join(codes, ","))
		}
		return nil
	},
}
