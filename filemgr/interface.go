package filemgr

import (
	"github.com/filecoin-project/venus-wallet/common"
	"github.com/filecoin-project/venus-wallet/config"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// file system
type Repo interface {
	// APIEndpoint returns multiaddress for communication with venus wallet API
	APIEndpoint() (multiaddr.Multiaddr, error)

	// APIToken returns JWT API Token for use in operations that require auth
	APIToken() ([]byte, error)

	APISecret() (*common.APIAlg, error)

	// APIStrategyToken cli pwd convert root token
	APIStrategyToken(password string) (string, error)

	Config() *config.Config

	AppendSupportAccount(newAccount string) error
}
