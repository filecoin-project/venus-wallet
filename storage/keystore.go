package storage

import (
	"fmt"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/crypto"
)

var (
	ErrKeyInfoNotFound = fmt.Errorf("key info not found")
	ErrKeyExists       = fmt.Errorf("key already exists")
)

// Constraint database implementation
// has: sqlite
type KeyStore interface {
	Put(key crypto.PrivateKey) error
	Get(addr core.Address) (crypto.PrivateKey, error)
	Has(addr core.Address) (bool, error)
	List() ([]core.Address, error)
	Delete(addr core.Address) error
}
