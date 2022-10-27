package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/filecoin-project/venus-wallet/config"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/crypto"
	"github.com/filecoin-project/venus-wallet/crypto/aes"
	"github.com/filecoin-project/venus-wallet/errcode"
	"github.com/filecoin-project/venus/venus-shared/api/permission"
	wallet_api "github.com/filecoin-project/venus/venus-shared/api/wallet"
	"github.com/filecoin-project/venus/venus-shared/types"
	"github.com/google/uuid"
)

var (
	ErrLocked          = errors.New("wallet locked")
	ErrPasswordEmpty   = errors.New("password not set")
	ErrInvalidPassword = errors.New("password mismatch")
	ErrPasswordExist   = errors.New("the password already exists")
	ErrAlreadyUnlocked = errors.New("wallet already unlocked")
	ErrAlreadyLocked   = errors.New("wallet already locked")
)

var EmptyPassword []byte

type DecryptFunc func(keyJson []byte, keyType types.KeyType) (crypto.PrivateKey, error)

// KeyMiddleware the middleware bridging strategy and wallet
type KeyMiddleware interface {
	// Encrypt aes encrypt key
	Encrypt(password []byte, key crypto.PrivateKey) (*aes.EncryptedKey, error)
	// Decrypt decrypt aes key
	Decrypt(password []byte, key *aes.EncryptedKey) (crypto.PrivateKey, error)
	// Next Check the password has been set and the wallet is locked
	Next() error
	// EqualRootToken compare the root token
	EqualRootToken(token string) error
	// CheckToken check if the `strategy` token has all permissions
	CheckToken(ctx context.Context) error
	wallet_api.IWalletLock
}

type KeyMixLayer struct {
	m         sync.RWMutex
	rootToken string // gen from password
	locked    bool
	password  []byte
	scryptN   int // aes cryptographic variable
	scryptP   int // aes cryptographic variable
}

func NewKeyMiddleware(cnf *config.CryptoFactor) KeyMiddleware {
	return &KeyMixLayer{
		locked:   true,
		password: nil,
		scryptN:  cnf.ScryptN,
		scryptP:  cnf.ScryptP,
	}
}

func (o *KeyMixLayer) SetPassword(ctx context.Context, password string) error {
	o.m.Lock()
	defer o.m.Unlock()
	if len(o.password) != 0 {
		return ErrPasswordExist
	}
	rootToken, err := o.genRootToken(ctx, password)
	if err != nil {
		return err
	}
	o.rootToken = rootToken
	hashPasswd := aes.Keccak256([]byte(password))
	o.password = hashPasswd
	o.locked = false
	return nil
}

func (o *KeyMixLayer) genRootToken(ctx context.Context, password string) (string, error) {
	hashPasswd := aes.Keccak256([]byte(password))
	rootKey, err := aes.EncryptData(hashPasswd, []byte("root"), o.scryptN, o.scryptP)
	if err != nil {
		return "", errors.New("failed to gen token seed")
	}
	rootKB, err := json.Marshal(rootKey)
	if err != nil {
		return "", errors.New("failed to marshal token seed")
	}
	rootk, err := uuid.NewRandomFromReader(bytes.NewBuffer(rootKB))
	if err != nil {
		return "", errors.New("failed to convert token seed to uuid")
	}
	return rootk.String(), nil
}

func (o *KeyMixLayer) EqualRootToken(token string) error {
	if len(o.password) == 0 || len(o.rootToken) == 0 {
		return ErrPasswordEmpty
	}
	if o.rootToken == token {
		return nil
	}
	return errcode.ErrWithoutPermission
}

func (o *KeyMixLayer) Unlock(ctx context.Context, password string) error {
	err := o.changeLock(password, false)
	if err != nil && err == ErrPasswordEmpty {
		return o.SetPassword(ctx, password)
	}
	return err
}

func (o *KeyMixLayer) Lock(ctx context.Context, password string) error {
	return o.changeLock(password, true)
}

func (o *KeyMixLayer) LockState(ctx context.Context) bool {
	return o.locked
}

func (o *KeyMixLayer) changeLock(password string, lock bool) error {
	o.m.Lock()
	defer o.m.Unlock()
	if len(o.password) == 0 {
		return ErrPasswordEmpty
	}
	if o.locked == lock {
		if o.locked {
			return ErrAlreadyLocked
		} else {
			return ErrAlreadyUnlocked
		}
	}
	hashPasswd := aes.Keccak256([]byte(password))
	if !bytes.Equal(o.password, hashPasswd) {
		return ErrInvalidPassword
	}
	o.locked = lock
	if !o.locked {
		o.password = hashPasswd
	}
	return nil
}

func (o *KeyMixLayer) CheckToken(ctx context.Context) error {
	if len(o.password) == 0 || len(o.rootToken) == 0 {
		return ErrPasswordEmpty
	}

	if core.WalletStrategyLevel == core.SLDisable || auth.HasPerm(ctx, permission.AllPermissions, permission.PermAdmin) {
		return nil
	}

	token := core.ContextStrategyToken(ctx)

	if o.rootToken == token {
		return nil
	}
	return errcode.ErrWithoutPermission
}

// Next Check the password has been set and the wallet is locked
func (o *KeyMixLayer) Next() error {
	o.m.RLock()
	defer o.m.RUnlock()
	if len(o.password) == 0 {
		return ErrPasswordEmpty
	}
	if o.locked {
		return ErrLocked
	}
	return nil
}

func (o *KeyMixLayer) Encrypt(password []byte, key crypto.PrivateKey) (*aes.EncryptedKey, error) {
	if len(password) == 0 {
		password = o.password
	}
	// EncryptKey encrypts a key using the specified scrypt parameters into a json
	// blob that can be decrypted later on.
	cryptoStruct, err := o.encryptData(password, key.Bytes())
	if err != nil {
		return nil, err
	}
	addr, _ := key.Address()
	encryptedKeyJSON := &aes.EncryptedKey{
		Address: addr.String(),
		KeyType: key.KeyType(),
		Crypto:  cryptoStruct,
	}
	return encryptedKeyJSON, nil
}

func (o *KeyMixLayer) encryptData(password []byte, data []byte) (*aes.CryptoJSON, error) {
	return aes.EncryptData(password, data, o.scryptN, o.scryptP)
}

func (o *KeyMixLayer) Decrypt(password []byte, key *aes.EncryptedKey) (crypto.PrivateKey, error) {
	if len(password) == 0 {
		password = o.password
	}
	// Depending on the version try to parse one way or another
	keyBytes, err := aes.Decrypt(key.Crypto, password)
	// Handle any decryption errors and return the key
	if err != nil {
		return nil, err
	}
	pkey, err := crypto.NewKeyFromData2(key.KeyType, keyBytes)
	if err != nil {
		return nil, err
	}
	return pkey, nil
}

func (o *KeyMixLayer) VerifyPassword(_ context.Context, password string) error {
	if bytes.Equal(o.password, aes.Keccak256([]byte(password))) {
		return nil
	}
	return ErrInvalidPassword
}
