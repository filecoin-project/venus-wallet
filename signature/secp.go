package signature

import (
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-crypto"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/minio/blake2b-simd"
)

func secpVerify(sig []byte, a core.Address, msg []byte) error {
	b2sum := blake2b.Sum256(msg)
	pubk, err := crypto.EcRecover(b2sum[:], sig)
	if err != nil {
		return err
	}

	maybeaddr, err := address.NewSecp256k1Address(pubk)
	if err != nil {
		return err
	}

	if a != maybeaddr {
		return fmt.Errorf("signature did not match")
	}

	return nil
}
