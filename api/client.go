package api

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"net/http"
)

func NewCommonRPC(ctx context.Context, addr string, requestHeader http.Header) (ICommon, jsonrpc.ClientCloser, error) {
	var res CommonAuth
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{
			&res.internal,
		},
		requestHeader,
	)

	return &res, closer, err
}

// NewFullNodeRPC creates a new http jsonrpc client.
func NewFullNodeRPC(ctx context.Context, addr string, requestHeader http.Header) (IFullAPI, jsonrpc.ClientCloser, error) {
	var res ServerAuth
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{
			&res.CommonAuth.internal,
			&res.WalletAuth.internal,
		}, requestHeader)

	return &res, closer, err
}

// nolint
func NewWalletRPC(ctx context.Context, addr string, requestHeader http.Header) (storage.IWallet, jsonrpc.ClientCloser, error) {
	var res WalletAuth
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{
			&res.internal,
		},
		requestHeader,
	)
	return &res, closer, err
}
