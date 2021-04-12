package cli_strategy

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/ipfs-force-community/venus-wallet/cli/helper"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/errcode"
	"github.com/urfave/cli/v2"
	"strconv"
	"strings"
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
		for k, v := range core.MethodNameList {
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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()
		tmp, err := api.GetMsgTypeTemplate(ctx, name)
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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)

		name := cctx.Args().First()
		tmp, err := api.GetMethodTemplateByName(ctx, name)
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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)

		name := cctx.Args().First()
		tmp, err := api.GetKeyBindByName(ctx, name)
		if err != nil {
			return err
		}
		var codes []string
		linq.From(core.FindCode(tmp.MetaTypes)).SelectT(func(i int) string {
			return strconv.FormatInt(int64(i), 10)
		}).ToSlice(&codes)
		fmt.Printf("address\t: %s\n", tmp.Address)
		fmt.Printf("types\t: %s\n", strings.Join(codes, ","))
		fmt.Printf("methods\t: %s\n", strings.Join(tmp.Methods, ","))
		return nil
	},
}

var strategyGetKeyBinds = &cli.Command{
	Name:      "keyBinds",
	Aliases:   []string{"kbs"},
	Usage:     "show keyBinds by address",
	ArgsUsage: "[address]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)

		address := cctx.Args().First()
		arr, err := api.GetKeyBinds(ctx, address)
		if err != nil {
			return err
		}
		for k, v := range arr {
			var codes []string
			linq.From(core.FindCode(v.MetaTypes)).SelectT(func(i int) string {
				return strconv.FormatInt(int64(i), 10)
			}).ToSlice(&codes)
			fmt.Printf("num\t: %d\n", k+1)
			fmt.Printf("name\t: %s\n", v.Name)
			fmt.Printf("addr\t: %s\n", v.Address)
			fmt.Printf("types\t: %s\n", strings.Join(codes, ","))
			fmt.Printf("methods\t: %s\n\n", strings.Join(v.Methods, ","))
		}
		return nil
	},
}

var strategyGetGroupByName = &cli.Command{
	Name:      "group",
	Aliases:   []string{"g"},
	Usage:     "show group by name",
	ArgsUsage: "[name]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)

		name := cctx.Args().First()
		group, err := api.GetGroupByName(ctx, name)
		if err != nil {
			return err
		}
		for k, v := range group.KeyBinds {
			var codes []string
			linq.From(core.FindCode(v.MetaTypes)).SelectT(func(i int) string {
				return strconv.FormatInt(int64(i), 10)
			}).ToSlice(&codes)
			fmt.Printf("num\t: %d\n", k+1)
			fmt.Printf("keybind\t: %s\n", v.Name)
			fmt.Printf("types\t: %s\n", strings.Join(codes, ","))
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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)

		from, to, err := helper.ReqFromTo(cctx, 0)
		if err != nil {
			return err
		}
		groups, err := api.ListGroups(ctx, int(from), int(to))
		if err != nil {
			return err
		}
		for k, v := range groups {
			fmt.Printf("%d\t: %s\n", k+1, v.Name)
		}
		return nil
	},
}

var strategyGroupTokens = &cli.Command{
	Name:      "groupTokens",
	Aliases:   []string{"gts"},
	Usage:     "show a range of tokens belong to group",
	ArgsUsage: "[groupName]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		groupName := cctx.Args().First()

		api, closer, err := helper.GetFullAPIWithPWD(cctx)
		if err != nil {
			return err
		}
		ctx := helper.ReqContext(cctx)
		defer closer()

		tks, err := api.GetWalletTokensByGroup(ctx, groupName)
		if err != nil {
			return err
		}
		for _, v := range tks {
			fmt.Println(v)
		}
		return nil
	},
}

var strategyTokenInfo = &cli.Command{
	Name:      "tokenInfo",
	Aliases:   []string{"ti"},
	Usage:     "show info about token",
	ArgsUsage: "[token]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		token := cctx.Args().First()

		api, closer, err := helper.GetFullAPIWithPWD(cctx)
		if err != nil {
			return err
		}
		ctx := helper.ReqContext(cctx)
		defer closer()

		ti, err := api.GetWalletTokenInfo(ctx, token)
		if err != nil {
			return err
		}
		fmt.Printf("groupName: %s\n", ti.Name)
		fmt.Println("keyBinds:")
		for k, v := range ti.KeyBinds {
			var codes []string
			linq.From(core.FindCode(v.MetaTypes)).SelectT(func(i int) string {
				return strconv.FormatInt(int64(i), 10)
			}).ToSlice(&codes)
			fmt.Printf("\tnum\t: %d\n", k+1)
			fmt.Printf("\tname\t: %s\n", v.Name)
			fmt.Printf("\taddr\t: %s\n", v.Address)
			fmt.Printf("\ttypes\t: %s\n", strings.Join(codes, ","))
			fmt.Printf("\tmethods\t: %s\n\n", strings.Join(v.Methods, ","))
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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)

		from, to, err := helper.ReqFromTo(cctx, 0)
		if err != nil {
			return err
		}
		arr, err := api.ListKeyBinds(ctx, int(from), int(to))
		if err != nil {
			return err
		}
		for k, v := range arr {
			var codes []string
			linq.From(core.FindCode(v.MetaTypes)).SelectT(func(i int) string {
				return strconv.FormatInt(int64(i), 10)
			}).ToSlice(&codes)
			fmt.Printf("num\t: %d\n", k+1)
			fmt.Printf("name\t: %s\n", v.Name)
			fmt.Printf("addr\t: %s\n", v.Address)
			fmt.Printf("types\t: %s\n", strings.Join(codes, ","))
			fmt.Printf("methods\t: %s\n\n", strings.Join(v.Methods, ","))
		}
		return nil
	},
}

var strategyListMethodTemplates = &cli.Command{
	Name:      "listMethodTemplates",
	Aliases:   []string{"lmt"},
	Usage:     "show a range of method templates",
	ArgsUsage: "[from to]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)

		from, to, err := helper.ReqFromTo(cctx, 0)
		if err != nil {
			return err
		}
		arr, err := api.ListMethodTemplates(ctx, int(from), int(to))
		if err != nil {
			return err
		}
		for k, v := range arr {
			fmt.Printf("num\t: %d\n", k+1)
			fmt.Printf("name\t: %s\n", v.Name)
			fmt.Printf("methods\t: %s\n\n", strings.Join(v.Methods, ","))
		}
		return nil
	},
}
var strategyListMsgTypeTemplates = &cli.Command{
	Name:      "listMsgTypeTemplates",
	Aliases:   []string{"lmtt"},
	Usage:     "show a range of method templates",
	ArgsUsage: "[from to]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)

		from, to, err := helper.ReqFromTo(cctx, 0)
		if err != nil {
			return err
		}
		arr, err := api.ListMsgTypeTemplates(ctx, from, to)
		if err != nil {
			return err
		}
		for k, v := range arr {
			var codes []string
			linq.From(core.FindCode(v.MetaTypes)).SelectT(func(i int) string {
				return strconv.FormatInt(int64(i), 10)
			}).ToSlice(&codes)
			fmt.Printf("num\t: %d\n", k+1)
			fmt.Printf("name\t: %s\n", v.Name)
			fmt.Printf("types\t: %s\n\n", strings.Join(codes, ","))
		}
		return nil
	},
}
