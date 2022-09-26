package filemgr

import (
	"github.com/filecoin-project/venus-wallet/config"
	"github.com/gbrlsnchs/jwt/v3"
)

// file system
type Repo interface {
	// APIEndpoint returns multiaddress for communication with venus wallet API
	APIEndpoint() (string, error)

	// APIToken returns JWT API Token for use in operations that require auth
	APIToken() ([]byte, error)

	APISecret() (*jwt.HMACSHA, error)

	// APIStrategyToken cli pwd convert root token
	APIStrategyToken(password string) (string, error)

	Config() *config.Config

	AppendSupportAccount(newAccount string) error
}
