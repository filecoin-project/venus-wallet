package crypto

import (
	"fmt"
	"github.com/ipfs-force-community/venus-wallet/core"
)

type PrivateKey interface {
	Public() []byte
	Sign([]byte) (*core.Signature, error)
	Bytes() []byte
	Address() (core.Address, error)
	Type() core.SigType
	KeyType() core.KeyType
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
