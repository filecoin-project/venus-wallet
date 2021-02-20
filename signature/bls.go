package signature

import (
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/crypto"
	"golang.org/x/xerrors"
)

// sig []byte, sigGroupcheck bool, pk []byte, pkValidate bool, msg Message, dst []byte,
func blsVerify(sig []byte, a core.Address, msg []byte) error {
	if !new(crypto.Signature).VerifyCompressed(sig, false, a.Payload()[:], false, msg, []byte(crypto.DST)) {
		return xerrors.New("bls signature failed to verify")
	}
	return nil
}
