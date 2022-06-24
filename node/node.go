package node

import (
	"context"
	"errors"

	"github.com/asaskevich/EventBus"
	"github.com/filecoin-project/venus-wallet/api/remotecli/httpparse"
	"github.com/filecoin-project/venus-wallet/config"
	v1 "github.com/filecoin-project/venus/venus-shared/api/chain/v1"
	"github.com/filecoin-project/venus/venus-shared/utils"

	"github.com/filecoin-project/venus-wallet/core"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("node")

// NodeClient connect Lotus or Venus node and call json RPC API
type NodeClient struct {
	// NOTE:
	Full   v1.FullNode
	Cancel func()
}

var EmptyNodeClient = &NodeClient{}

func NewNodeClient(cnf *config.StrategyConfig) (*NodeClient, error) {
	if cnf.Level < core.SLMethod {
		return EmptyNodeClient, nil
	}
	if cnf.NodeURL == core.StringEmpty {
		return nil, errors.New("nod/ip4e url can not be empty when level is SLMethod")
	}
	ai, err := httpparse.ParseApiInfo(cnf.NodeURL)
	if err != nil {
		return nil, err
	}
	addr, err := ai.DialArgs()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	full, closer, err := v1.NewFullNodeRPC(ctx, addr, ai.AuthHeader())
	if err != nil {
		return nil, err
	}
	cli := &NodeClient{Full: full, Cancel: closer}

	log.Info("node client initialize successfully")

	if err = reloadMethodNames(ctx, full); err != nil {
		return nil, err
	}

	return cli, nil
}

func reloadMethodNames(ctx context.Context, full v1.FullNode) error {
	if err := utils.LoadBuiltinActors(ctx, full); err != nil {
		return err
	}
	core.ReloadMethodNames()

	return nil
}

func NewEventBus() EventBus.Bus {
	return EventBus.New()
}
