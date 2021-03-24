package cli_strategy

import (
	"github.com/urfave/cli/v2"
)

var StrategyCmd = &cli.Command{
	Name:    "strategy",
	Usage:   "Manage logging",
	Aliases: []string{"st"},
	Subcommands: []*cli.Command{
		//read
		strategyTypeList,
		strategyMethodList,
		strategyGetMsgTypeTemplate,
		strategyGetMethodTemplateByName,
		strategyGetKeyBindByName,
		strategyGetKeyBinds,
		strategyGetGroupByName,
		strategyListGroup,
		strategyGroupTokens,
		strategyListKeyBinds,
		strategyTokenInfo,
		strategyListMethodTemplates,
		strategyListMsgTypeTemplates,
		//write
		strategyNewMsgTypeTemplate,
		strategyNewMethodTemplate,
		strategyNewKeyBindCustom,
		strategyNewKeyBindFromTemplate,
		strategyNewGroup,
		strategyNewWalletToken,
		strategyRemoveMsgTypeTemplate,
		strategyRemoveMethodTemplate,
		strategyRemoveKeyBind,
		strategyRemoveKeyBindByAddress,
		strategyRemoveGroup,
		strategyRemoveToken,
		strategyPullMethodIntoKeyBind,
		strategyPullMsgTypeFromKeyBind,
		strategyPushMethodIntoKeyBind,
		strategyPushMsgTypeIntoKeyBind,
	},
}
