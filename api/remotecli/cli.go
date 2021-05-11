package remotecli

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/venus-wallet/api"
	"github.com/filecoin-project/venus-wallet/common"
	"github.com/filecoin-project/venus-wallet/storage/wallet"
	"net/http"
)

// NewWalletRPC RPCClient returns an RPC client connected to a node
// @addr			reference ./httpparse/ParseApiInfo()
// @requestHeader 	reference ./httpparse/ParseApiInfo()
func NewWalletRPC(ctx context.Context, addr string, requestHeader http.Header) (wallet.IWallet, jsonrpc.ClientCloser, error) {
	var res api.WalletAPIAdapter
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{
			&res.Internal,
		},
		requestHeader,
	)
	return &res, closer, err
}

func NewCommonRPC(ctx context.Context, addr string, requestHeader http.Header) (common.ICommon, jsonrpc.ClientCloser, error) {
	var res api.CommonAPIAdapter
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{
			&res.Internal,
		},
		requestHeader,
	)

	return &res, closer, err
}

// NewFullNodeRPC creates a new httpparse jsonrpc remotecli.
func NewFullNodeRPC(ctx context.Context, addr string, requestHeader http.Header) (api.IFullAPI, jsonrpc.ClientCloser, error) {
	var res api.FullAPIAdapter
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{
			&res.CommonAPIAdapter.Internal,
			&res.WalletAPIAdapter.Internal,
			&res.WalletLockAPIAdapter.Internal,
			&res.StrategyAPIAdapter.Internal,
		}, requestHeader)

	return &res, closer, err
}
