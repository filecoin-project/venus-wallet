package crypto

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-crypto"
	"github.com/filecoin-project/venus/venus-shared/types"
	"github.com/minio/blake2b-simd"

	c1 "github.com/filecoin-project/go-state-types/crypto"
)

type secpPrivateKey struct {
	key []byte
}

func newSecpKeyFromData(data []byte) PrivateKey {
	return &secpPrivateKey{
		key: data,
	}
}

func genSecpPrivateKey() (PrivateKey, error) {
	prv, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	p := &secpPrivateKey{
		key: prv,
	}
	return p, nil
}

func (p *secpPrivateKey) Public() []byte {
	return crypto.PublicKey(p.key)
}

func (p *secpPrivateKey) Sign(msg []byte) (*c1.Signature, error) {
	b2sum := blake2b.Sum256(msg)
	sig, err := crypto.Sign(p.key, b2sum[:])
	if err != nil {
		return nil, err
	}
	return &c1.Signature{
		Data: sig,
		Type: p.Type(),
	}, nil
}
func (p *secpPrivateKey) Bytes() []byte {
	return p.key
}
func (p *secpPrivateKey) Address() (address.Address, error) {
	addr, err := address.NewSecp256k1Address(p.Public())
	if err != nil {
		return address.Undef, fmt.Errorf("converting Secp256k1 to address: %w", err)
	}
	return addr, nil
}
func (p *secpPrivateKey) Type() types.SigType {
	return types.SigTypeSecp256k1
}
func (p *secpPrivateKey) KeyType() types.KeyType {
	return types.KTSecp256k1
}
func (p *secpPrivateKey) ToKeyInfo() *types.KeyInfo {
	return &types.KeyInfo{
		PrivateKey: p.Bytes(),
		Type:       types.KTSecp256k1,
	}
}
func secpVerify(sig []byte, a address.Address, msg []byte) error {
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
