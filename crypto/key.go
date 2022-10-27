package crypto

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/venus/venus-shared/types"
)

// private key constraints in wallet
type PrivateKey interface {
	// Private to public
	Public() []byte
	// private key signature
	Sign([]byte) (*crypto.Signature, error)
	// private key data
	Bytes() []byte
	// key address, depends on network changes
	Address() (address.Address, error)
	// key sign type
	Type() types.SigType
	// key type
	KeyType() types.KeyType
	// map to keyInfo
	ToKeyInfo() *types.KeyInfo
}

func Verify(sig *crypto.Signature, addr address.Address, msg []byte) error {
	if sig == nil {
		return fmt.Errorf("signature is nil")
	}
	if addr.Protocol() == address.ID {
		return fmt.Errorf("must resolve ID addresses before using them to verify a signature")
	}
	switch sig.Type {
	case types.SigTypeSecp256k1:
		return secpVerify(sig.Data, addr, msg)
	case types.SigTypeBLS:
		return blsVerify(sig.Data, addr, msg)
	default:
		return fmt.Errorf("cannot verify signature of unsupported type: %v", sig.Type)
	}
}

func GeneratePrivateKey(st types.SigType) (PrivateKey, error) {
	switch st {
	case types.SigTypeSecp256k1:
		return genSecpPrivateKey()
	case types.SigTypeBLS:
		return genBlsPrivate()
	default:
		return nil, fmt.Errorf("invalid signature type: %d", st)
	}
}

func NewKeyFromKeyInfo(ki *types.KeyInfo) (PrivateKey, error) {
	switch ki.Type {
	case types.KTBLS:
		return newBlsKeyFromData(ki.PrivateKey)
	case types.KTSecp256k1:
		return newSecpKeyFromData(ki.PrivateKey), nil
	default:
		return nil, fmt.Errorf("invalid key type: %s", ki.Type)
	}
}

func NewKeyFromData2(kt types.KeyType, prv []byte) (PrivateKey, error) {
	switch kt {
	case types.KTBLS:
		return newBlsKeyFromData(prv)
	case types.KTSecp256k1:
		return newSecpKeyFromData(prv), nil
	default:
		return nil, fmt.Errorf("invalid key type: %s", kt)
	}
}

func NewKeyFromData(st types.SigType, prv []byte) (PrivateKey, error) {
	switch st {
	case types.SigTypeSecp256k1:
		return newSecpKeyFromData(prv), nil
	case types.SigTypeBLS:
		return newBlsKeyFromData(prv)
	default:
		return nil, fmt.Errorf("invalid signature type: %d", st)
	}
}
