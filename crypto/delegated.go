package crypto

import (
	"fmt"

	"golang.org/x/crypto/sha3"

	"github.com/filecoin-project/go-address"
	gocrypto "github.com/filecoin-project/go-crypto"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/venus/venus-shared/types"
)

type delegatedPrivateKey struct {
	key []byte
}

func newDelegatedKeyFromData(data []byte) PrivateKey {
	return &delegatedPrivateKey{
		key: data,
	}
}

func genDelegatedPrivateKey() (PrivateKey, error) {
	prv, err := gocrypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	p := &delegatedPrivateKey{
		key: prv,
	}
	return p, nil
}

func (p *delegatedPrivateKey) Public() []byte {
	return gocrypto.PublicKey(p.key)
}

func (p *delegatedPrivateKey) Sign(msg []byte) (*crypto.Signature, error) {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(msg)
	hashSum := hasher.Sum(nil)
	sig, err := gocrypto.Sign(p.key, hashSum)
	if err != nil {
		return nil, err
	}

	return &crypto.Signature{
		Type: p.Type(),
		Data: sig,
	}, nil
}

func (p *delegatedPrivateKey) Bytes() []byte {
	return p.key
}

func (p *delegatedPrivateKey) Address() (address.Address, error) {
	pubKey := p.Public()
	// Transitory Delegated signature verification as per FIP-0055
	ethAddr, err := types.EthAddressFromPubKey(pubKey)
	if err != nil {
		return address.Undef, fmt.Errorf("failed to calculate Eth address from public key: %w", err)
	}
	ea, err := types.CastEthAddress(ethAddr)
	if err != nil {
		return address.Undef, fmt.Errorf("failed to create ethereum address from bytes: %w", err)
	}

	return ea.ToFilecoinAddress()
}

func (p *delegatedPrivateKey) Type() types.SigType {
	return types.SigTypeDelegated
}

func (p *delegatedPrivateKey) KeyType() types.KeyType {
	return types.KTDelegated
}

func (p *delegatedPrivateKey) ToKeyInfo() *types.KeyInfo {
	return &types.KeyInfo{
		PrivateKey: p.Bytes(),
		Type:       types.KTSecp256k1,
	}
}

func delegatedVerify(sig []byte, a address.Address, msg []byte) error {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(msg)
	hash := hasher.Sum(nil)

	pubk, err := gocrypto.EcRecover(hash, sig)
	if err != nil {
		return err
	}

	// if we get an uncompressed public key (that's what we get from the library,
	// but putting this check here for defensiveness), strip the prefix
	if pubk[0] == 0x04 {
		pubk = pubk[1:]
	}

	hasher.Reset()
	hasher.Write(pubk)
	addrHash := hasher.Sum(nil)

	// The address hash will not start with [12]byte{0xff}, so we don't have to use
	// EthAddr.ToFilecoinAddress() to handle the case with an id address
	// Also, importing ethtypes here will cause circulating import
	maybeaddr, err := address.NewDelegatedAddress(builtin.EthereumAddressManagerActorID, addrHash[12:])
	if err != nil {
		return err
	}

	if maybeaddr != a {
		return fmt.Errorf("signature did not match")
	}

	return nil
}
