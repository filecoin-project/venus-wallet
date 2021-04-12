package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ipfs-force-community/venus-wallet/core"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
	"hash"
	"io"
)

var (
	ErrDecrypt = errors.New("could not decrypt key with given password")
)

const (
	keyHeaderKDF = "scrypt"
	// memory and taking approximately 1s CPU time on a modern processor.
	StandardScryptN = 1 << 18

	// StandardScryptP is the P parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	StandardScryptP = 1

	scryptR     = 8
	scryptDKLen = 32
)

type cipherparamsJSON struct {
	IV string `json:"iv"`
}

type CryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams cipherparamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type EncryptedKey struct {
	Address string       `json:"address"`
	KeyType core.KeyType `json:"type"`
	Crypto  *CryptoJSON  `json:"crypto"`
}

func getKDFKey(cryptoJSON *CryptoJSON, auth []byte) ([]byte, error) {
	salt, err := hex.DecodeString(cryptoJSON.KDFParams["salt"].(string))
	if err != nil {
		return nil, err
	}
	dkLen := ensureInt(cryptoJSON.KDFParams["dklen"])

	if cryptoJSON.KDF == keyHeaderKDF {
		n := ensureInt(cryptoJSON.KDFParams["n"])
		r := ensureInt(cryptoJSON.KDFParams["r"])
		p := ensureInt(cryptoJSON.KDFParams["p"])
		return scrypt.Key(auth, salt, n, r, p, dkLen)

	} else if cryptoJSON.KDF == "pbkdf2" {
		c := ensureInt(cryptoJSON.KDFParams["c"])
		prf := cryptoJSON.KDFParams["prf"].(string)
		if prf != "hmac-sha256" {
			return nil, fmt.Errorf("unsupported PBKDF2 PRF: %s", prf)
		}
		key := pbkdf2.Key(auth, salt, c, dkLen, sha256.New)
		return key, nil
	}

	return nil, fmt.Errorf("unsupported KDF: %s", cryptoJSON.KDF)
}
func ensureInt(x interface{}) int {
	res, ok := x.(int)
	if !ok {
		res = int(x.(float64))
	}
	return res
}

func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	// AES-128 is selected due to size of encryptKey.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}

// nolint
func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := NewKeccakState()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}

// NewKeccakState creates a new KeccakState
func NewKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}

type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

func EncryptData(password, data []byte, scryptN, scryptP int) (*CryptoJSON, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	derivedKey, err := scrypt.Key(password, salt, scryptN, scryptR, scryptP, scryptDKLen)
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
	mac := Keccak256(derivedKey[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptN
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

func Decrypt(cryptoJson *CryptoJSON, pwd []byte) ([]byte, error) {
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

	derivedKey, err := getKDFKey(cryptoJson, pwd)
	if err != nil {
		return nil, err
	}

	calculatedMAC := Keccak256(derivedKey[16:32], cipherText)
	if !bytes.Equal(calculatedMAC, mac) {
		return nil, ErrDecrypt
	}

	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, err
	}
	return plainText, err
}
