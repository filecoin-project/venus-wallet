package httpparse

import (
	"github.com/ipfs-force-community/venus-common-utils/apiinfo"
	"golang.org/x/xerrors"
	"net/http"
	"strings"
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
	if len(sep) == 3 {
		return &APIInfo{
			Addr:  sep[1],
			Token: []byte(sep[0] + ":" + sep[2]),
		}, nil
	} else if len(sep) == 2 {
		return &APIInfo{
			Addr:  sep[1],
			Token: []byte(sep[0]),
		}, nil
	} else {
		return nil, xerrors.Errorf("invalidate api info string %s", s)
	}

}

func (a APIInfo) DialArgs() (string, error) {
	return apiinfo.DialArgs(a.Addr, "v0")
}

func (a APIInfo) AuthHeader() http.Header {
	if len(a.Token) != 0 {
		headers := http.Header{}
		headers.Add(ServiceToken, "Bearer "+string(a.Token))
		return headers
	}
	return nil
}
