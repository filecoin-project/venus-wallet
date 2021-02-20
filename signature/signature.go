package signature

import (
	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/venus-wallet/core"
	"golang.org/x/xerrors"
)

//nolint
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
