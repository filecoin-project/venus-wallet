package api

import (
	"encoding/binary"
	"github.com/filecoin-project/go-state-types/cbor"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/lotus/api"
	"github.com/minio/blake2b-simd"
	"golang.org/x/xerrors"
	fcbor "github.com/fxamacker/cbor"
	"io"
	"io/ioutil"
)

type FullNode interface {
	Common
	api.WalletAPI
}

const (
	MTDeals = api.MsgType("deals")
	// extra is nil, 'toSign' is cbor raw bytes of 'DrawRandomParams'
	//  following types follow above rule
	MTDrawRandomParam   = api.MsgType("drawrandomparam")
	MTSignedVoucher     = api.MsgType("signedvoucher")
	MTStorageAsk        = api.MsgType("storageask")
	MTAskResponse       = api.MsgType("askresponse")
	MTNetWorkResponse   = api.MsgType("networkresposne")
	MTProviderDealState = api.MsgType("providerdealstate")

	// reference : storagemarket/impl/client.go:330
	// sign storagemarket.ClientDeal.ProposalCid,
	// MsgMeta.Extra is nil, 'toSign' is market.ClientDealProposal
	// storagemarket.ClientDeal.ProposalCid equals cborutil.AsIpld(market.ClientDealProposal).Cid()
	MTClientDeal = api.MsgType("clientdeal")
)

type DrawRandomParams struct {
	Rbase   []byte
	Pers    crypto.DomainSeparationTag
	Round   abi.ChainEpoch
	Entropy []byte
}

// return store.DrawRandomness(dr.Rbase, dr.Pers, dr.Round, dr.Entropy)
func (dr *DrawRandomParams) SignBytes() ([]byte, error) {
	h := blake2b.New256()
	if err := binary.Write(h, binary.BigEndian, int64(dr.Pers)); err != nil {
		return nil, xerrors.Errorf("deriving randomness: %w", err)
	}
	VRFDigest := blake2b.Sum256(dr.Rbase)
	_, err := h.Write(VRFDigest[:])
	if err != nil {
		return nil, xerrors.Errorf("hashing VRFDigest: %w", err)
	}
	if err := binary.Write(h, binary.BigEndian, dr.Round); err != nil {
		return nil, xerrors.Errorf("deriving randomness: %w", err)
	}
	_, err = h.Write(dr.Entropy)
	if err != nil {
		return nil, xerrors.Errorf("hashing entropy: %w", err)
	}

	return h.Sum(nil), nil
}

func (dr *DrawRandomParams) MarshalCBOR(w io.Writer) error {
	data, err := fcbor.Marshal(dr, fcbor.EncOptions{})
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (dr *DrawRandomParams) UnmarshalCBOR(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return fcbor.Unmarshal(data, dr)
}

var _ = cbor.Unmarshaler((*DrawRandomParams)(nil))
var _ = cbor.Marshaler((*DrawRandomParams)(nil))
