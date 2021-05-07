package cli

import (
	"github.com/filecoin-project/venus-wallet/cli/cli_strategy"
	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	authCmd,
	logCmd,
	cli_strategy.StrategyCmd,
	walletNew,
	walletList,
	walletExport,
	walletImport,
	walletSign,
	walletDel,
	walletSetPassword,
	walletUnlock,
	walletlock,
	walletLockState,
}
