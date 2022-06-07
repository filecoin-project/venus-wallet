package crypto

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus/venus-shared/types"
	blst "github.com/supranational/blst/bindings/go"
)

const DST = string("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_NUL_")

type SecretKey = blst.SecretKey
type PublicKey = blst.P1Affine
type Signature = blst.P2Affine
type AggregateSignature = blst.P2Aggregate

type blsPrivate struct {
	key    *SecretKey
	public []byte
}

func newBlsKeyFromData(data []byte) (PrivateKey, error) {
	pk := new(SecretKey).FromLEndian(data)
	if pk == nil || !pk.Valid() {
		return nil, errors.New("bls signature invalid private key")
	}
	return &blsPrivate{
		key:    pk,
		public: new(PublicKey).From(pk).Compress(),
	}, nil
}

func genBlsPrivate() (PrivateKey, error) {
	// Generate 32 bytes of randomness
	var ikm [32]byte
	_, err := rand.Read(ikm[:])
	if err != nil {
		return nil, fmt.Errorf("bls signature error generating random data")
	}
	// Note private keys seem to be serialized little-endian!
	sk := blst.KeyGen(ikm[:])
	pk := &blsPrivate{
		key:    sk,
		public: new(PublicKey).From(sk).Compress(),
	}

	return pk, nil
}

func (p *blsPrivate) Public() []byte {
	return p.public
}

func (p *blsPrivate) Sign(msg []byte) (*core.Signature, error) {
	return &core.Signature{
		Data: new(Signature).Sign(p.key, msg, []byte(DST)).Compress(),
		Type: p.Type(),
	}, nil
}
func (p *blsPrivate) Bytes() []byte {
	return p.key.ToLEndian()
}
func (p *blsPrivate) Address() (core.Address, error) {
	addr, err := address.NewBLSAddress(p.public)
	if err != nil {
		return core.NilAddress, fmt.Errorf("converting BLS to address: %w", err)
	}
	return addr, nil
}

func (p *blsPrivate) Type() types.SigType {
	return types.SigTypeBLS
}
func (p *blsPrivate) KeyType() types.KeyType {
	return types.KTBLS
}
func (p *blsPrivate) ToKeyInfo() *types.KeyInfo {
	return &types.KeyInfo{
		PrivateKey: p.Bytes(),
		Type:       types.KTBLS,
	}
}

// sig []byte, sigGroupcheck bool, pk []byte, pkValidate bool, msg Message, dst []byte,
func blsVerify(sig []byte, a core.Address, msg []byte) error {
	if !new(Signature).VerifyCompressed(sig, false, a.Payload()[:], false, msg, []byte(DST)) {
		return errors.New("bls signature failed to verify")
	}
	return nil
}
