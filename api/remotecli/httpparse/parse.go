package httpparse

import (
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"net/http"
	"regexp"
	"strings"
)

var (
	infoWithToken = regexp.MustCompile("^[a-zA-Z0-9\\-_]+?\\.[a-zA-Z0-9\\-_]+?\\.([a-zA-Z0-9\\-_]+)?:.+$")                    //nolint
	strategyToken = regexp.MustCompile("[0-9a-fA-F]{8}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{4}\\-[0-9a-fA-F]{12}") //nolint
)

const (
	ServiceToken        = "Authorization"
	WalletStrategyToken = "StrategyToken"
)

type APIInfo struct {
	Addr          multiaddr.Multiaddr
	Token         []byte
	StrategyToken []byte
}

func ParseApiInfo(s string) (*APIInfo, error) {

	var (
		tok        []byte
		strategyTk []byte
	)

	if infoWithToken.Match([]byte(s)) {
		sp := strings.SplitN(s, ":", 2)
		tok = []byte(sp[0])
		s = sp[1]
	}
	if strategyToken.Match([]byte(s)) {
		sp := strings.SplitN(s, ":", 2)
		strategyTk = []byte(sp[1])
		s = sp[0]
	}
	strma := strings.TrimSpace(s)
	apima, err := multiaddr.NewMultiaddr(strma)
	if err != nil {
		return nil, err
	}
	return &APIInfo{
		Addr:          apima,
		Token:         tok,
		StrategyToken: strategyTk,
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
		headers.Add(ServiceToken, "Bearer "+string(a.Token))
		headers.Add(WalletStrategyToken, string(a.StrategyToken))
		return headers
	}
	return nil
}
