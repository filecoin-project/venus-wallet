package api

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
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
			&res.internal,
			&res.wallet,
		}, requestHeader)

	return &res, closer, err
}
