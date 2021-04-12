package cli_strategy

import (
	"errors"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/ipfs-force-community/venus-wallet/cli/helper"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/errcode"
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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		groupName := cctx.Args().First()

		api, closer, err := helper.GetFullAPIWithPWD(cctx)
		if err != nil {
			return err
		}
		ctx := helper.ReqContext(cctx)
		defer closer()

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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
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
	ArgsUsage: "[name, address, codes, <methods>]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 3 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()
		address := cctx.Args().Get(1)
		codesStr := cctx.Args().Get(2)
		var codes []int
		for _, v := range strings.Split(codesStr, ",") {
			code, err := strconv.Atoi(v)
			if err != nil {
				return errors.New("codes must be int type")
			}
			codes = append(codes, code)
		}
		methods := make([]string, 0)
		if cctx.NArg() == 4 {
			methodsStr := cctx.Args().Get(3)
			methods = strings.Split(methodsStr, ",")
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
	ArgsUsage: "[name, address, msgTypeTemplateName, <methodTemplateName>]",
	Action: func(cctx *cli.Context) error {
		if cctx.NArg() < 3 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		name := cctx.Args().First()
		address := cctx.Args().Get(1)
		mttName := cctx.Args().Get(2)
		mtName := ""
		if cctx.NArg() == 4 {
			mtName = cctx.Args().Get(3)
		}

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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		api, closer, err := helper.GetFullAPI(cctx)
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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		name := cctx.Args().First()

		api, closer, err := helper.GetFullAPIWithPWD(cctx)
		if err != nil {
			return err
		}
		ctx := helper.ReqContext(cctx)
		defer closer()

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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		name := cctx.Args().First()

		api, closer, err := helper.GetFullAPIWithPWD(cctx)
		if err != nil {
			return err
		}
		ctx := helper.ReqContext(cctx)
		defer closer()

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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		name := cctx.Args().First()

		api, closer, err := helper.GetFullAPIWithPWD(cctx)
		if err != nil {
			return err
		}
		ctx := helper.ReqContext(cctx)
		defer closer()

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
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		token := cctx.Args().First()

		api, closer, err := helper.GetFullAPIWithPWD(cctx)
		if err != nil {
			return err
		}
		ctx := helper.ReqContext(cctx)
		defer closer()

		err = api.RemoveToken(ctx, token)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyPullMsgTypeFromKeyBind = &cli.Command{
	Name:      "pullMsgTypeFromKeyBind",
	Aliases:   []string{"pullMT4KB"},
	Usage:     "remove elements of msgTypes in keyBind",
	ArgsUsage: "[keyBindName, code1 code2 ...]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		name := cctx.Args().First()
		codes := make([]int, 0)
		for _, arg := range cctx.Args().Slice()[1:] {
			code, err := strconv.Atoi(arg)
			if err != nil {
				return errors.New("code must be the number")
			}
			codes = append(codes, code)
		}

		api, closer, err := helper.GetFullAPIWithPWD(cctx)
		if err != nil {
			return err
		}
		ctx := helper.ReqContext(cctx)
		defer closer()

		kb, err := api.PullMsgTypeFromKeyBind(ctx, name, codes)
		if err != nil {
			return err
		}
		var codesOut []string
		linq.From(core.FindCode(kb.MetaTypes)).SelectT(func(i int) string {
			return strconv.FormatInt(int64(i), 10)
		}).ToSlice(&codesOut)
		fmt.Printf("address\t: %s\n", kb.Address)
		fmt.Printf("types\t: %s\n", strings.Join(codesOut, ","))
		fmt.Printf("methods\t: %s\n", strings.Join(kb.Methods, ","))
		return nil
	},
}

var strategyPullMethodIntoKeyBind = &cli.Command{
	Name:      "pullMethodFromKeyBind",
	Aliases:   []string{"pullM4KB"},
	Usage:     "remove elements of methods in keyBind",
	ArgsUsage: "[keyBindName, method1 method2 ...]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		name := cctx.Args().First()
		methods := make([]string, 0)
		methods = append(methods, cctx.Args().Slice()[1:]...)

		api, closer, err := helper.GetFullAPIWithPWD(cctx)
		if err != nil {
			return err
		}
		ctx := helper.ReqContext(cctx)
		defer closer()

		kb, err := api.PullMethodFromKeyBind(ctx, name, methods)
		if err != nil {
			return err
		}
		var codesOut []string
		linq.From(core.FindCode(kb.MetaTypes)).SelectT(func(i int) string {
			return strconv.FormatInt(int64(i), 10)
		}).ToSlice(&codesOut)
		fmt.Printf("address\t: %s\n", kb.Address)
		fmt.Printf("types\t: %s\n", strings.Join(codesOut, ","))
		fmt.Printf("methods\t: %s\n", strings.Join(kb.Methods, ","))
		return nil
	},
}

var strategyPushMsgTypeIntoKeyBind = &cli.Command{
	Name:      "pushMsgTypeIntoKeyBind",
	Aliases:   []string{"pushMT2KB"},
	Usage:     "append msgTypes into keyBind",
	ArgsUsage: "[keyBindName, code1 code2 ...]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		name := cctx.Args().First()
		codes := make([]int, 0)
		for _, arg := range cctx.Args().Slice()[1:] {
			code, err := strconv.Atoi(arg)
			if err != nil {
				return errors.New("code must be the number")
			}
			codes = append(codes, code)
		}

		api, closer, err := helper.GetFullAPIWithPWD(cctx)
		if err != nil {
			return err
		}
		ctx := helper.ReqContext(cctx)
		defer closer()

		kb, err := api.PushMsgTypeIntoKeyBind(ctx, name, codes)
		if err != nil {
			return err
		}
		var codesOut []string
		linq.From(core.FindCode(kb.MetaTypes)).SelectT(func(i int) string {
			return strconv.FormatInt(int64(i), 10)
		}).ToSlice(&codesOut)
		fmt.Printf("address\t: %s\n", kb.Address)
		fmt.Printf("types\t: %s\n", strings.Join(codesOut, ","))
		fmt.Printf("methods\t: %s\n", strings.Join(kb.Methods, ","))
		return nil
	},
}

var strategyPushMethodIntoKeyBind = &cli.Command{
	Name:      "pushMethodIntoKeyBind",
	Aliases:   []string{"pushM2KB"},
	Usage:     "append methods into keyBind",
	ArgsUsage: "[keyBindName, method1 method2 ...]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		name := cctx.Args().First()
		methods := make([]string, 0)
		methods = append(methods, cctx.Args().Slice()[1:]...)

		api, closer, err := helper.GetFullAPIWithPWD(cctx)
		if err != nil {
			return err
		}
		ctx := helper.ReqContext(cctx)
		defer closer()

		kb, err := api.PushMethodIntoKeyBind(ctx, name, methods)
		if err != nil {
			return err
		}
		var codesOut []string
		linq.From(core.FindCode(kb.MetaTypes)).SelectT(func(i int) string {
			return strconv.FormatInt(int64(i), 10)
		}).ToSlice(&codesOut)
		fmt.Printf("address\t: %s\n", kb.Address)
		fmt.Printf("types\t: %s\n", strings.Join(codesOut, ","))
		fmt.Printf("methods\t: %s\n", strings.Join(kb.Methods, ","))
		return nil
	},
}
