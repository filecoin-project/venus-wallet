package impl

import (
	logging "github.com/ipfs/go-log/v2"

	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/node/impl/common"
	"github.com/ipfs-force-community/venus-wallet/node/impl/full"
)

var log = logging.Logger("node")

type FullNodeAPI struct {
	common.CommonAPI
	full.WalletAPI
}

var _ api.FullNode = &FullNodeAPI{}
