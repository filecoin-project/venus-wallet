package cli_strategy

import (
	"github.com/urfave/cli/v2"
)

var StrategyCmd = &cli.Command{
	Name:  "strategy",
	Usage: "Manage logging",
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
		strategyListKeyBinds,
		strategyListMethodTemplates,
		strategyListMsgTypeTemplates,
		//write
		strategyNewMsgTypeTemplate,
		strategyNewMethodTemplate,
		strategyNewKeyBindCustom,
		strategyNewKeyBindFromTemplate,
		strategyNewGroup,
		strategyRemoveMsgTypeTemplate,
		strategyRemoveMethodTemplate,
		strategyRemoveKeyBind,
		strategyRemoveKeyBindByAddress,
		strategyRemoveGroup,
	},
}
