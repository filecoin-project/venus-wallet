package httpparse

import (
	"net/http"
	"strings"

	"github.com/filecoin-project/venus/venus-shared/api"
	"golang.org/x/xerrors"
)

const (
	ServiceToken = "Authorization"
)

// APIInfo parse URL string to
type APIInfo struct {
	Addr  string
	Token []byte
}

func ParseApiInfo(s string) (*APIInfo, error) {
	sep := strings.Split(s, ":")
	if len(sep) != 2 {
		return nil, xerrors.Errorf("invalidate api info string %s", s)
	}
	return &APIInfo{
		Addr:  sep[1],
		Token: []byte(sep[0]),
	}, nil
}

func (a APIInfo) DialArgs() (string, error) {
	return api.DialArgs(a.Addr, "v0")
}

func (a APIInfo) AuthHeader() http.Header {
	if len(a.Token) != 0 {
		headers := http.Header{}
		headers.Add(ServiceToken, "Bearer "+string(a.Token))
		return headers
	}
	return nil
}
