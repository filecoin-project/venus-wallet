package remotecli

import (
	"context"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
	apiutil "github.com/filecoin-project/venus/venus-shared/api"
	api "github.com/filecoin-project/venus/venus-shared/api/wallet"
)

// NewWalletRPC RPCClient returns an RPC client connected to a node
// @addr			reference ./httpparse/ParseApiInfo()
// @requestHeader 	reference ./httpparse/ParseApiInfo()
func NewWalletRPC(ctx context.Context, addr string, requestHeader http.Header) (api.IWallet, jsonrpc.ClientCloser, error) {
	var res api.IWalletStruct
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{
			&res.Internal,
		},
		requestHeader,
	)
	return &res, closer, err
}

func NewCommonRPC(ctx context.Context, addr string, requestHeader http.Header) (api.ICommon, jsonrpc.ClientCloser, error) {
	var res api.ICommonStruct
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
	var res api.IFullAPIStruct
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin", apiutil.GetInternalStructs(&res), requestHeader)

	return &res, closer, err
}
