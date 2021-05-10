package crypto

import (
	"fmt"
	"github.com/filecoin-project/venus-wallet/core"
)

// private key constraints in wallet
type PrivateKey interface {
	// Private to public
	Public() []byte
	// private key signature
	Sign([]byte) (*core.Signature, error)
	// private key data
	Bytes() []byte
	// key address, depends on network changes
	Address() (core.Address, error)
	// key sign type
	Type() core.SigType
	// key type
	KeyType() core.KeyType
	// map to keyInfo
	ToKeyInfo() *core.KeyInfo
}

func GeneratePrivateKey(st core.SigType) (PrivateKey, error) {
	switch st {
	case core.SigTypeSecp256k1:
		return genSecpPrivateKey()
	case core.SigTypeBLS:
		return genBlsPrivate()
	default:
		return nil, fmt.Errorf("invalid signature type: %d", st)
	}
}
func NewKeyFromKeyInfo(ki *core.KeyInfo) (PrivateKey, error) {
	switch ki.Type {
	case core.KTBLS:
		return newBlsKeyFromData(ki.PrivateKey)
	case core.KTSecp256k1:
		return newSecpKeyFromData(ki.PrivateKey), nil
	default:
		return nil, fmt.Errorf("invalid key type: %s", ki.Type)
	}
}

func NewKeyFromData2(kt core.KeyType, prv []byte) (PrivateKey, error) {
	switch kt {
	case core.KTBLS:
		return newBlsKeyFromData(prv)
	case core.KTSecp256k1:
		return newSecpKeyFromData(prv), nil
	default:
		return nil, fmt.Errorf("invalid key type: %s", kt)
	}
}

func NewKeyFromData(st core.SigType, prv []byte) (PrivateKey, error) {
	switch st {
	case core.SigTypeSecp256k1:
		return newSecpKeyFromData(prv), nil
	case core.SigTypeBLS:
		return newBlsKeyFromData(prv)
	default:
		return nil, fmt.Errorf("invalid signature type: %d", st)
	}
}
