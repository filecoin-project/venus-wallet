package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/filecoin-project/venus-wallet/config"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/crypto"
	"github.com/filecoin-project/venus-wallet/crypto/aes"
	"github.com/filecoin-project/venus-wallet/errcode"
	"github.com/google/uuid"
	"sync"
)

type IWalletLock interface {
	// SetPassword do it first after program setup
	SetPassword(ctx context.Context, password string) error
	// unlock the wallet and enable IWallet logic
	Unlock(ctx context.Context, password string) error
	// lock the wallet and disable IWallet logic
	Lock(ctx context.Context, password string) error
	// show lock state
	LockState(ctx context.Context) bool
}

var (
	ErrLocked          = errors.New("wallet locked")
	ErrPasswordEmpty   = errors.New("password not set")
	ErrInvalidPassword = errors.New("password mismatch")
	ErrPasswordExist   = errors.New("the password already exists")
	ErrAlreadyUnlocked = errors.New("wallet already unlocked")
	ErrAlreadyLocked   = errors.New("wallet already locked")
)

type DecryptFunc func(keyJson []byte, keyType core.KeyType) (crypto.PrivateKey, error)

// KeyMiddleware the middleware bridging strategy and wallet
type KeyMiddleware interface {
	Encrypt(key crypto.PrivateKey) (*aes.EncryptedKey, error)
	Decrypt(key *aes.EncryptedKey) (crypto.PrivateKey, error)
	Next() error
	EqualRootToken(token string) error
	CheckToken(ctx context.Context) error
	IWalletLock
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
		return core.StringEmpty, errors.New("failed to gen token seed")
	}
	rootKB, err := json.Marshal(rootKey)
	if err != nil {
		return core.StringEmpty, errors.New("failed to marshal token seed")
	}
	rootk, err := uuid.NewRandomFromReader(bytes.NewBuffer(rootKB))
	if err != nil {
		return core.StringEmpty, errors.New("failed to convert token seed to uuid")
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
	return o.changeLock(password, false)
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
	return nil
}

func (o *KeyMixLayer) CheckToken(ctx context.Context) error {
	token := core.ContextStrategyToken(ctx)
	if len(o.password) == 0 || len(o.rootToken) == 0 {
		return ErrPasswordEmpty
	}
	if core.WalletStrategyLevel == core.SLDisable {
		return nil
	}
	if o.rootToken == token {
		return nil
	}
	return errcode.ErrWithoutPermission
}

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

func (o *KeyMixLayer) Encrypt(key crypto.PrivateKey) (*aes.EncryptedKey, error) {
	// EncryptKey encrypts a key using the specified scrypt parameters into a json
	// blob that can be decrypted later on.
	cryptoStruct, err := o.encryptData(key.Bytes())
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

func (o *KeyMixLayer) encryptData(data []byte) (*aes.CryptoJSON, error) {
	return aes.EncryptData(o.password, data, o.scryptN, o.scryptP)
}

func (o *KeyMixLayer) Decrypt(key *aes.EncryptedKey) (crypto.PrivateKey, error) {
	// Depending on the version try to parse one way or another
	keyBytes, err := aes.Decrypt(key.Crypto, o.password)
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
