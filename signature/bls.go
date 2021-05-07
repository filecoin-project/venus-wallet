package signature

import (
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/crypto"
	"golang.org/x/xerrors"
)

// sig []byte, sigGroupcheck bool, pk []byte, pkValidate bool, msg Message, dst []byte,
func blsVerify(sig []byte, a core.Address, msg []byte) error {
	if !new(crypto.Signature).VerifyCompressed(sig, false, a.Payload()[:], false, msg, []byte(crypto.DST)) {
		return xerrors.New("bls signature failed to verify")
	}
	return nil
}
