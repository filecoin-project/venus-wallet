package storage

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-wallet/crypto/aes"
	"github.com/filecoin-project/venus/venus-shared/types"
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
	Get(addr address.Address) (*aes.EncryptedKey, error)
	// Has check the PrivateKey exist in the KeyStore
	Has(addr address.Address) (bool, error)
	// List lists all the keys stored in the KeyStore
	List() ([]address.Address, error)
	// Delete removes a key from keystore
	Delete(addr address.Address) error
}

type QueryParams = types.QuerySignRecordParams

type SignRecord = types.SignRecord

type IRecorder interface {
	Record(rcd *SignRecord) error
	QueryRecord(params *QueryParams) ([]SignRecord, error)
}
