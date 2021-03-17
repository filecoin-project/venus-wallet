package cli_strategy

import (
	"errors"
	"fmt"
	"github.com/ipfs-force-community/venus-wallet/cli/helper"
	"github.com/urfave/cli/v2"
	"strconv"
	"strings"
)

var strategyNewWalletToken = &cli.Command{
	Name:      "newWalletToken",
	Aliases:   []string{"newWT"},
	Usage:     "create a wallet token with group",
	ArgsUsage: "[groupName]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		groupName := cctx.Args().First()

		token, err := api.NewWalletToken(ctx, groupName)
		if err != nil {
			return err
		}
		fmt.Println(token)
		return nil
	},
}

var strategyNewMsgTypeTemplate = &cli.Command{
	Name:      "newMsgTypeTemplate",
	Aliases:   []string{"newMTT"},
	Usage:     "create a msgType common template",
	ArgsUsage: "[name, code1 code2 ...]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()
		codes := make([]int, 0)
		for _, arg := range cctx.Args().Slice()[1:] {
			code, err := strconv.Atoi(arg)
			if err != nil {
				return errors.New("code must be the number")
			}
			codes = append(codes, code)
		}
		err = api.NewMsgTypeTemplate(ctx, name, codes)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyNewMethodTemplate = &cli.Command{
	Name:      "newMethodTemplate",
	Aliases:   []string{"newMT"},
	Usage:     "create a msg methods common template",
	ArgsUsage: "[name, method1 method2 ...]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()
		methods := make([]string, 0)
		methods = append(methods, cctx.Args().Slice()[1:]...)
		err = api.NewMethodTemplate(ctx, name, methods)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyNewKeyBindCustom = &cli.Command{
	Name:      "newKeyBindCustom",
	Aliases:   []string{"newKBC"},
	Usage:     "create a strategy about wallet bind msgType and methods",
	ArgsUsage: "[name, address, codes, methods]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 4 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()
		address := cctx.Args().Get(1)
		codesStr := cctx.Args().Get(2)
		methodsStr := cctx.Args().Get(3)
		methods := strings.Split(methodsStr, ",")
		var codes []int
		for _, v := range strings.Split(codesStr, ",") {
			code, err := strconv.Atoi(v)
			if err != nil {
				return errors.New("codes must be int type")
			}
			codes = append(codes, code)
		}
		err = api.NewKeyBindCustom(ctx, name, address, codes, methods)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyNewKeyBindFromTemplate = &cli.Command{
	Name:      "newKeyBindFromTemplate",
	Aliases:   []string{"newKBFT"},
	Usage:     "create a strategy about wallet bind msgType and methods with template",
	ArgsUsage: "[name, address, msgTypeTemplateName, methodTemplateName]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 4 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()
		address := cctx.Args().Get(1)
		mttName := cctx.Args().Get(2)
		mtName := cctx.Args().Get(3)

		err = api.NewKeyBindFromTemplate(ctx, name, address, mttName, mtName)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyNewGroup = &cli.Command{
	Name:      "newGroup",
	Aliases:   []string{"newG"},
	Usage:     "create a group with keyBinds",
	ArgsUsage: "[name, keyBindName1 keyBindName2 ...]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()
		names := make([]string, 0)
		names = append(names, cctx.Args().Slice()[1:]...)
		err = api.NewGroup(ctx, name, names)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyRemoveMsgTypeTemplate = &cli.Command{
	Name:      "removeMsgTypeTemplate",
	Aliases:   []string{"rmMTT"},
	Usage:     "remove msgTypeTemplate ( not affect the group strategy that has been created)",
	ArgsUsage: "[name]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()
		err = api.RemoveMsgTypeTemplate(ctx, name)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyRemoveMethodTemplate = &cli.Command{
	Name:      "removeMethodTemplate",
	Aliases:   []string{"rmMT"},
	Usage:     "remove MethodTemplate ( not affect the group strategy that has been created)",
	ArgsUsage: "[name]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()
		err = api.RemoveMethodTemplate(ctx, name)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyRemoveKeyBind = &cli.Command{
	Name:      "removeKeyBind",
	Aliases:   []string{"rmKB"},
	Usage:     "remove keyBind ( not affect the group strategy that has been created)",
	ArgsUsage: "[name]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()
		err = api.RemoveKeyBind(ctx, name)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyRemoveKeyBindByAddress = &cli.Command{
	Name:      "removeKeyBindByAddress",
	Aliases:   []string{"rmKBBA"},
	Usage:     "remove keyBinds by address ( not affect the group strategy that has been created)",
	ArgsUsage: "[name]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()

		num, err := api.RemoveKeyBindByAddress(ctx, name)
		if err != nil {
			return err
		}
		fmt.Printf("%d rows of data were deleted", num)
		return nil
	},
}

var strategyRemoveGroup = &cli.Command{
	Name:      "removeGroup",
	Aliases:   []string{"rmG"},
	Usage:     "remove group by address ( not affect the group strategy that has been created)",
	ArgsUsage: "[name]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()

		err = api.RemoveGroup(ctx, name)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyRemoveToken = &cli.Command{
	Name:      "removeToken",
	Aliases:   []string{"rmT"},
	Usage:     "remove token",
	ArgsUsage: "[token]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 1 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		token := cctx.Args().First()

		err = api.RemoveToken(ctx, token)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}
