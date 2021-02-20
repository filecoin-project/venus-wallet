package filemgr

import (
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/config"
	"github.com/multiformats/go-multiaddr"
)

type Repo interface {
	// APIEndpoint returns multiaddress for communication with venus wallet API
	APIEndpoint() (multiaddr.Multiaddr, error)

	// APIToken returns JWT API Token for use in operations that require auth
	APIToken() ([]byte, error)

	APISecret() (*api.APIAlg, error)

	Config() *config.Config
}
