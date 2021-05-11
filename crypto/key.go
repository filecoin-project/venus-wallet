package crypto

import (
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-wallet/core"
	"golang.org/x/xerrors"
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

func Verify(sig *core.Signature, addr core.Address, msg []byte) error {
	if sig == nil {
		return xerrors.Errorf("signature is nil")
	}
	if addr.Protocol() == address.ID {
		return xerrors.Errorf("must resolve ID addresses before using them to verify a signature")
	}
	switch sig.Type {
	case core.SigTypeSecp256k1:
		return secpVerify(sig.Data, addr, msg)
	case core.SigTypeBLS:
		return blsVerify(sig.Data, addr, msg)
	default:
		return xerrors.Errorf("cannot verify signature of unsupported type: %v", sig.Type)
	}
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
