package storage

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/config"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/crypto"
	"golang.org/x/crypto/scrypt"
	"io"
	"sync"
)

var (
	ErrLocked          = errors.New("wallet locked")
	ErrPasswordEmpty   = errors.New("password not set")
	ErrInvalidPassword = errors.New("password mismatch")
	ErrPasswordExist   = errors.New("the password already exists")
)

type DecryptFunc func(keyJson []byte, keyType core.KeyType) (crypto.PrivateKey, error)

type KeyMiddleware interface {
	Encrypt(key crypto.PrivateKey) (*EncryptedKey, error)
	Decrypt(key *EncryptedKey) (crypto.PrivateKey, error)
	Next() error
	api.IWalletLock
}

type KeyMixLayer struct {
	m        sync.RWMutex
	cache    map[core.Address]crypto.PrivateKey
	locked   bool
	password []byte
	scryptN  int
	scryptP  int
}

func NewKeyMiddleware(cnf *config.CryptoFactor) KeyMiddleware {
	return &KeyMixLayer{
		cache:    make(map[core.Address]crypto.PrivateKey),
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
	hashPasswd := keccak256([]byte(password))
	o.password = hashPasswd
	o.locked = false
	return nil
}

func (o *KeyMixLayer) Unlock(ctx context.Context, password string) error {
	return o.changeLock(password, false)
}
func (o *KeyMixLayer) Lock(ctx context.Context, password string) error {
	return o.changeLock(password, true)
}
func (o *KeyMixLayer) changeLock(password string, lock bool) error {
	o.m.Lock()
	defer o.m.Unlock()
	if len(o.password) == 0 {
		return ErrPasswordEmpty
	}
	if o.locked == lock {
		return nil
	}
	hashPasswd := keccak256([]byte(password))
	if !bytes.Equal(o.password, hashPasswd) {
		return ErrInvalidPassword
	}
	o.locked = lock
	return nil
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

func (o *KeyMixLayer) Encrypt(key crypto.PrivateKey) (*EncryptedKey, error) {
	// EncryptKey encrypts a key using the specified scrypt parameters into a json
	// blob that can be decrypted later on.
	cryptoStruct, err := o.encryptData(key.Bytes())
	if err != nil {
		return nil, err
	}
	addr, _ := key.Address()
	encryptedKeyJSON := &EncryptedKey{
		Address: addr.String(),
		KeyType: key.KeyType(),
		Crypto:  cryptoStruct,
	}
	return encryptedKeyJSON, nil
}

// Encryptdata encrypts the data given as 'data' with the password 'auth'.
func (o *KeyMixLayer) encryptData(data []byte) (*CryptoJSON, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	derivedKey, err := scrypt.Key(o.password, salt, o.scryptN, scryptR, o.scryptP, scryptDKLen)
	if err != nil {
		return nil, err
	}
	encryptKey := derivedKey[:16]

	iv := make([]byte, aes.BlockSize) // 16
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	cipherText, err := aesCTRXOR(encryptKey, data, iv)
	if err != nil {
		return nil, err
	}
	mac := keccak256(derivedKey[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = o.scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = o.scryptP
	scryptParamsJSON["dklen"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)
	cipherParamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}
	cryptoStruct := &CryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}
	return cryptoStruct, nil
}
func (o *KeyMixLayer) Decrypt(key *EncryptedKey) (crypto.PrivateKey, error) {
	// Depending on the version try to parse one way or another
	keyBytes, err := o.decrypt(key.Crypto)
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

func (o *KeyMixLayer) decrypt(cryptoJson *CryptoJSON) ([]byte, error) {
	if cryptoJson.Cipher != "aes-128-ctr" {
		return nil, fmt.Errorf("cipher not supported: %v", cryptoJson.Cipher)
	}
	mac, err := hex.DecodeString(cryptoJson.MAC)
	if err != nil {
		return nil, err
	}

	iv, err := hex.DecodeString(cryptoJson.CipherParams.IV)
	if err != nil {
		return nil, err
	}

	cipherText, err := hex.DecodeString(cryptoJson.CipherText)
	if err != nil {
		return nil, err
	}

	derivedKey, err := getKDFKey(cryptoJson, o.password)
	if err != nil {
		return nil, err
	}

	calculatedMAC := keccak256(derivedKey[16:32], cipherText)
	if !bytes.Equal(calculatedMAC, mac) {
		return nil, ErrDecrypt
	}

	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, err
	}
	return plainText, err
}
