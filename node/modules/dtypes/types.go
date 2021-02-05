package dtypes

import (
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/ipfs/go-datastore"
	"github.com/multiformats/go-multiaddr"
)

// MetadataDS stores metadata
// dy default it's namespaced under /metadata in main repo datastore
type MetadataDS datastore.Batching
type NetworkName string

type APIAlg jwt.HMACSHA

type APIEndpoint multiaddr.Multiaddr
