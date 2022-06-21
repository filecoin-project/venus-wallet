package node

import (
	"context"
	"errors"
	"fmt"

	"github.com/asaskevich/EventBus"
	"github.com/filecoin-project/venus-wallet/api/remotecli/httpparse"
	"github.com/filecoin-project/venus-wallet/config"
	builtinactors "github.com/filecoin-project/venus/venus-shared/builtin-actors"
	types "github.com/filecoin-project/venus/venus-shared/types"
	"github.com/filecoin-project/venus/venus-shared/utils"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/venus-wallet/core"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("node")

// NodeClient connect Lotus or Venus node and call json RPC API
type NodeClient struct {
	// NOTE:
	StateGetActor    func(ctx context.Context, actor core.Address, tsk types.TipSetKey) (*types.Actor, error)
	StateNetworkName func(ctx context.Context) (types.NetworkName, error)
	Cancel           func()
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
	cli := &NodeClient{}
	ctx := context.Background()
	closer, err := jsonrpc.NewClient(ctx, addr, "Filecoin", cli, ai.AuthHeader())
	if err != nil {
		return nil, err
	}
	cli.Cancel = closer
	log.Info("node client initialize successfully")

	if err = reloadMethodNames(ctx, cli); err != nil {
		return nil, err
	}

	return cli, nil
}

func reloadMethodNames(ctx context.Context, cli *NodeClient) error {
	networkName, err := cli.StateNetworkName(ctx)
	if err != nil {
		return fmt.Errorf("failed to got network namee %v", err)
	}

	nt, err := utils.NetworkNameToNetworkType(networkName)
	if err != nil {
		if errors.Is(err, utils.ErrMay2kNetwork) {
			nt = types.Network2k
		} else {
			return err
		}
	}
	if err := builtinactors.SetNetworkBundle(nt); err != nil {
		return err
	}
	core.ReloadMethodNames()

	return nil
}

func NewEventBus() EventBus.Bus {
	return EventBus.New()
}
