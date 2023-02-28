package cli

import (
	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	authCmd,
	logCmd,
	walletNew,
	walletList,
	walletExport,
	walletImport,
	walletSign,
	walletDel,
	walletSetPassword,
	walletUnlock,
	walletLock,
	walletLockState,
	supportCmds,
	recordCmd,
}
