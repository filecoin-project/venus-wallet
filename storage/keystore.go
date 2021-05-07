package storage

import (
	"fmt"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/crypto/aes"
)

var (
	ErrKeyInfoNotFound = fmt.Errorf("key info not found")
	ErrKeyExists       = fmt.Errorf("key already exists")
)

// Constraint database implementation
// has: sqlite
type KeyStore interface {
	// Put saves a key info
	Put(key *aes.EncryptedKey) error
	// Get gets a key out of keystore and returns PrivateKey corresponding to key address
	Get(addr core.Address) (*aes.EncryptedKey, error)
	// Has check the PrivateKey exist in the KeyStore
	Has(addr core.Address) (bool, error)
	// List lists all the keys stored in the KeyStore
	List() ([]core.Address, error)
	// Delete removes a key from keystore
	Delete(addr core.Address) error
}
