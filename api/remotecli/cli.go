package remotecli

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/common"
	"github.com/ipfs-force-community/venus-wallet/storage/wallet"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"net/http"
	"regexp"
	"strings"
)

// RPCClient returns an RPC client connected to a node
func NewWalletRPC(ctx context.Context, addr string, requestHeader http.Header) (wallet.IWallet, jsonrpc.ClientCloser, error) {
	var res api.WalletAuth
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{
			&res.Internal,
		},
		requestHeader,
	)
	return &res, closer, err
}

func NewCommonRPC(ctx context.Context, addr string, requestHeader http.Header) (common.ICommon, jsonrpc.ClientCloser, error) {
	var res api.CommonAuth
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{
			&res.Internal,
		},
		requestHeader,
	)

	return &res, closer, err
}

// NewFullNodeRPC creates a new http jsonrpc remotecli.
func NewFullNodeRPC(ctx context.Context, addr string, requestHeader http.Header) (api.IFullAPI, jsonrpc.ClientCloser, error) {
	var res api.ServiceAuth
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{
			&res.CommonAuth.Internal,
			&res.WalletAuth.Internal,
			&res.WalletLockAuth.Internal,
		}, requestHeader)

	return &res, closer, err
}

var (
	infoWithToken = regexp.MustCompile("^[a-zA-Z0-9\\-_]+?\\.[a-zA-Z0-9\\-_]+?\\.([a-zA-Z0-9\\-_]+)?:.+$") //nolint
)

type APIInfo struct {
	Addr          multiaddr.Multiaddr
	Token         []byte
	StrategyToken []byte
}

func ParseApiInfo(s string) (*APIInfo, error) {
	var tok []byte
	if infoWithToken.Match([]byte(s)) {
		sp := strings.SplitN(s, ":", 2)
		tok = []byte(sp[0])
		s = sp[1]
	}
	strma := strings.TrimSpace(s)
	apima, err := multiaddr.NewMultiaddr(strma)
	if err != nil {
		return nil, err
	}
	return &APIInfo{
		Addr:  apima,
		Token: tok,
	}, nil
}

func (a APIInfo) DialArgs() (string, error) {
	_, addr, err := manet.DialArgs(a.Addr)
	if strings.HasPrefix(addr, "0.0.0.0:") {
		addr = "127.0.0.1:" + addr[8:]
	}
	return "ws://" + addr + "/rpc/v0", err
}

func (a APIInfo) AuthHeader() http.Header {
	if len(a.Token) != 0 {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer "+string(a.Token))
		headers.Add("StrategyToken", string(a.StrategyToken))
		return headers
	}
	return nil
}
