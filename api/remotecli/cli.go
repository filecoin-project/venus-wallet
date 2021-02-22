package remotecli

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"net/http"
	"regexp"
	"strings"
)

// nolint
func NewWalletRPC(ctx context.Context, addr string, requestHeader http.Header) (api.IWallet, jsonrpc.ClientCloser, error) {
	var res api.WalletAuth
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{
			&res.Internal,
		},
		requestHeader,
	)
	return &res, closer, err
}

var (
	infoWithToken = regexp.MustCompile("^[a-zA-Z0-9\\-_]+?\\.[a-zA-Z0-9\\-_]+?\\.([a-zA-Z0-9\\-_]+)?:.+$") //nolint
)

type APIInfo struct {
	Addr  multiaddr.Multiaddr
	Token []byte
}

// nolint
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
		return headers
	}
	return nil
}
