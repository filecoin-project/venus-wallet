package crypto

import (
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-crypto"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/minio/blake2b-simd"
	"golang.org/x/xerrors"
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

func (p *secpPrivateKey) Sign(msg []byte) (*core.Signature, error) {
	b2sum := blake2b.Sum256(msg)
	sig, err := crypto.Sign(p.key, b2sum[:])
	if err != nil {
		return nil, err
	}
	return &core.Signature{
		Data: sig,
		Type: p.Type(),
	}, nil
}
func (p *secpPrivateKey) Bytes() []byte {
	return p.key
}
func (p *secpPrivateKey) Address() (core.Address, error) {
	addr, err := address.NewSecp256k1Address(p.Public())
	if err != nil {
		return core.NilAddress, xerrors.Errorf("converting Secp256k1 to address: %w", err)
	}
	return addr, nil
}
func (p *secpPrivateKey) Type() core.SigType {
	return core.SigTypeSecp256k1
}
func (p *secpPrivateKey) KeyType() core.KeyType {
	return core.KTSecp256k1
}
func (p *secpPrivateKey) ToKeyInfo() *core.KeyInfo {
	return &core.KeyInfo{
		PrivateKey: p.Bytes(),
		Type:       core.KTSecp256k1,
	}
}
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
