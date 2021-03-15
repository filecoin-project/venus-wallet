package cli_strategy

import (
	"errors"
	"fmt"
	"github.com/ipfs-force-community/venus-wallet/cli/helper"
	"github.com/urfave/cli/v2"
	"strconv"
	"strings"
)

var strategyNewMsgTypeTemplate = &cli.Command{
	Name:      "newMsgTypeTemplate",
	Aliases:   []string{"newMTT"},
	Usage:     "create a msgType common template",
	ArgsUsage: "[name, code1 code2 ...]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
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
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		err = api.NewMsgTypeTemplate(name, codes)
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
		name := cctx.Args().First()
		methods := make([]string, 0)
		methods = append(methods, cctx.Args().Slice()[1:]...)
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		err = api.NewMethodTemplate(name, methods)
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
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		err = api.NewKeyBindCustom(name, address, codes, methods)
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
		name := cctx.Args().First()
		address := cctx.Args().Get(1)
		mttName := cctx.Args().Get(2)
		mtName := cctx.Args().Get(3)

		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		err = api.NewKeyBindFromTemplate(name, address, mttName, mtName)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyNewGroup = &cli.Command{
	Name:      "NewGroup",
	Aliases:   []string{"newG"},
	Usage:     "create a group with keyBinds",
	ArgsUsage: "[name, keyBindName1 keyBindName2 ...]",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() < 2 {
			return helper.ShowHelp(cctx, ErrParameterMismatch)
		}
		name := cctx.Args().First()
		names := make([]string, 0)
		names = append(names, cctx.Args().Slice()[1:]...)
		api, closer, err := helper.GetFullNodeAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		err = api.NewGroup(name, names)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyRemoveMsgTypeTemplate = &cli.Command{
	Name:      "RemoveMsgTypeTemplate",
	Aliases:   []string{"rmMTT"},
	Usage:     "remove msgTypeTemplate ( not affect the group strategy that has been created)",
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
		err = api.RemoveMsgTypeTemplate(name)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyRemoveMethodTemplate = &cli.Command{
	Name:      "RemoveMethodTemplate",
	Aliases:   []string{"rmMT"},
	Usage:     "remove MethodTemplate ( not affect the group strategy that has been created)",
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
		err = api.RemoveMethodTemplate(name)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyRemoveKeyBind = &cli.Command{
	Name:      "RemoveKeyBind",
	Aliases:   []string{"rmKB"},
	Usage:     "remove keyBind ( not affect the group strategy that has been created)",
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
		err = api.RemoveKeyBind(name)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}

var strategyRemoveKeyBindByAddress = &cli.Command{
	Name:      "RemoveKeyBindByAddress",
	Aliases:   []string{"rmKBBA"},
	Usage:     "remove keyBinds by address ( not affect the group strategy that has been created)",
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
		num, err := api.RemoveKeyBindByAddress(name)
		if err != nil {
			return err
		}
		fmt.Printf("%d rows of data were deleted", num)
		return nil
	},
}

var strategyRemoveGroup = &cli.Command{
	Name:      "RemoveGroup",
	Aliases:   []string{"rmG"},
	Usage:     "remove group by address ( not affect the group strategy that has been created)",
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
		err = api.RemoveGroup(name)
		if err != nil {
			return err
		}
		fmt.Println("success")
		return nil
	},
}
