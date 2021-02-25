package core

import (
	"bytes"
	"github.com/filecoin-project/go-state-types/abi"
	fbig "github.com/filecoin-project/go-state-types/big"
	proof2 "github.com/filecoin-project/specs-actors/v2/actors/runtime/proof"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	cbg "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"
	"io"
)

var (
	lengthBufBlock         = []byte{144}
	lengthBufTicket        = []byte{129}
	lengthBufElectionProof = []byte{130}
	lengthBufBeaconEntry   = []byte{130}
	blockHeaderCIDLen      int
)

func init() {
	// hash a large string of zeros so we don't estimate based on inlined CIDs.
	var buf [256]byte
	c, err := abi.CidBuilder.Sum(buf[:])
	if err != nil {
		panic(err)
	}
	blockHeaderCIDLen = len(c.Bytes())
}

type VRFPi []byte
type Ticket struct {
	// A proof output by running a VRF on the VRFProof of the parent ticket
	VRFProof VRFPi
}
type ElectionProof struct {
	WinCount int64

	// A proof output by running a VRF on the VRFProof of the parent ticket
	VRFProof VRFPi
}
type BeaconEntry struct {
	Round uint64
	Data  []byte
}

type TipSetKey struct {
	value string
}

// Block is a block in the blockchain.
type Block struct {
	// Miner is the address of the miner actor that mined this block.
	Miner Address `json:"miner"`

	// Ticket is the ticket submitted with this block.
	Ticket Ticket `json:"ticket"`

	// ElectionProof is the vrf proof giving this block's miner authoring rights
	ElectionProof *ElectionProof `json:"electionProof"`

	// BeaconEntries contain the verifiable oracle randomness used to elect
	// this block's author leader
	BeaconEntries []*BeaconEntry `json:"beaconEntries"`

	// WinPoStProof are the winning post proofs
	WinPoStProof []proof2.PoStProof `json:"winPoStProof"`

	// Parents is the set of parents this block was based on. Typically one,
	// but can be several in the case where there were multiple winning ticket-
	// holders for an epoch.
	Parents TipSetKey `json:"parents"`

	// ParentWeight is the aggregate chain weight of the parent set.
	ParentWeight fbig.Int `json:"parentWeight"`

	// Height is the chain height of this block.
	Height abi.ChainEpoch `json:"height"`

	// ParentStateRoot is the CID of the root of the state tree after application of the messages in the parent tipset
	// to the parent tipset's state root.
	ParentStateRoot Cid `json:"parentStateRoot,omitempty"`

	// ParentMessageReceipts is a list of receipts corresponding to the application of the messages in the parent tipset
	// to the parent tipset's state root (corresponding to this block's ParentStateRoot).
	ParentMessageReceipts Cid `json:"parentMessageReceipts,omitempty"`

	// Messages is the set of messages included in this block
	Messages Cid `json:"messages,omitempty"`

	// The aggregate signature of all BLS signed messages in the block
	BLSAggregate *Signature `json:"BLSAggregate"`

	// The timestamp, in seconds since the Unix epoch, at which this block was created.
	Timestamp uint64 `json:"timestamp"`

	// The signature of the miner's worker key over the block
	BlockSig *Signature `json:"blocksig"`

	// ForkSignaling is extra data used by miners to communicate
	ForkSignaling uint64 `json:"forkSignaling"`

	ParentBaseFee TokenAmount `json:"parentBaseFee"`
}

func (b *Block) SignatureData() []byte {
	tmp := &Block{
		Miner:                 b.Miner,
		Ticket:                b.Ticket,
		ElectionProof:         b.ElectionProof,
		Parents:               b.Parents,
		ParentWeight:          b.ParentWeight,
		Height:                b.Height,
		Messages:              b.Messages,
		ParentStateRoot:       b.ParentStateRoot,
		ParentMessageReceipts: b.ParentMessageReceipts,
		WinPoStProof:          b.WinPoStProof,
		BeaconEntries:         b.BeaconEntries,
		Timestamp:             b.Timestamp,
		BLSAggregate:          b.BLSAggregate,
		ForkSignaling:         b.ForkSignaling,
		ParentBaseFee:         b.ParentBaseFee,
		// BlockSig omitted
	}

	return tmp.rawData()
}

const BLAKE2B_MIN = 0xb201

// The multihash function identifier to use for content addresses.
const DefaultHashFunction = uint64(BLAKE2B_MIN + 31)

// A builder for all blockchain CIDs.
// Note that sector commitments use a different scheme.
var DefaultCidBuilder = cid.V1Builder{Codec: cid.DagCBOR, MhType: DefaultHashFunction}

func (b *Block) rawData() []byte {
	buf := new(bytes.Buffer)
	err := b.MarshalCBOR(buf)
	if err != nil {
		panic(err)
	}
	data := buf.Bytes()
	c, err := DefaultCidBuilder.Sum(data)
	if err != nil {
		panic(err)
	}

	blk, err := blocks.NewBlockWithCid(data, c)
	if err != nil {
		panic(err)
	}
	n, err := cbor.DecodeBlock(blk)
	if err != nil {
		panic(err)
	}
	return n.RawData()
}

func (t *Ticket) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufTicket); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.VRFProof (block.VRFPi) (slice)
	if len(t.VRFProof) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.VRFProof was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.VRFProof))); err != nil {
		return err
	}

	if _, err := w.Write(t.VRFProof[:]); err != nil {
		return err
	}
	return nil
}

func (t *ElectionProof) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufElectionProof); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.WinCount (int64) (int64)
	if t.WinCount >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.WinCount)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.WinCount-1)); err != nil {
			return err
		}
	}

	// t.VRFProof (block.VRFPi) (slice)
	if len(t.VRFProof) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.VRFProof was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.VRFProof))); err != nil {
		return err
	}

	if _, err := w.Write(t.VRFProof[:]); err != nil {
		return err
	}
	return nil
}

func (t *BeaconEntry) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufBeaconEntry); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Round (uint64) (uint64)

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.Round)); err != nil {
		return err
	}

	// t.Data ([]uint8) (slice)
	if len(t.Data) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Data was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.Data))); err != nil {
		return err
	}

	if _, err := w.Write(t.Data[:]); err != nil {
		return err
	}
	return nil
}

/*var lengthBufPoStProof = []byte{130}

func (t *PoStProof) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufPoStProof); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.PoStProof (abi.RegisteredPoStProof) (int64)
	if t.PoStProof >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.PoStProof)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.PoStProof-1)); err != nil {
			return err
		}
	}

	// t.ProofBytes ([]uint8) (slice)
	if len(t.ProofBytes) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.ProofBytes was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.ProofBytes))); err != nil {
		return err
	}

	if _, err := w.Write(t.ProofBytes[:]); err != nil {
		return err
	}
	return nil
}*/

func decodeKey(encoded []byte) ([]Cid, error) {
	// To avoid reallocation of the underlying array, estimate the number of CIDs to be extracted
	// by dividing the encoded length by the expected CID length.
	estimatedCount := len(encoded) / blockHeaderCIDLen
	cids := make([]Cid, 0, estimatedCount)
	nextIdx := 0
	for nextIdx < len(encoded) {
		nr, c, err := cid.CidFromBytes(encoded[nextIdx:])
		if err != nil {
			return nil, err
		}
		cids = append(cids, c)
		nextIdx += nr
	}
	return cids, nil
}

// Cids returns a slice of the CIDs comprising this key.
func (tipsetKey TipSetKey) Cids() []Cid {
	cids, err := decodeKey([]byte(tipsetKey.value))
	if err != nil {
		panic("invalid tipset key: " + err.Error())
	}
	return cids
}
func (tipsetKey TipSetKey) MarshalCBOR(w io.Writer) error {
	cids := tipsetKey.Cids()
	if len(cids) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Parents was too long")
	}
	scratch := make([]byte, 9)

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(cids))); err != nil {
		return err
	}
	for _, v := range cids {
		if err := cbg.WriteCidBuf(scratch, w, v); err != nil {
			return xerrors.Errorf("failed writing cid field t.Parents: %v", err)
		}
	}
	return nil
}

func (t *Block) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufBlock); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Miner (address.Address) (struct)
	if err := t.Miner.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Ticket (block.Ticket) (struct)
	if err := t.Ticket.MarshalCBOR(w); err != nil {
		return err
	}

	// t.ElectionProof (block.ElectionProof) (struct)
	if err := t.ElectionProof.MarshalCBOR(w); err != nil {
		return err
	}

	// t.BeaconEntries ([]*block.BeaconEntry) (slice)
	if len(t.BeaconEntries) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.BeaconEntries was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.BeaconEntries))); err != nil {
		return err
	}
	for _, v := range t.BeaconEntries {
		if err := v.MarshalCBOR(w); err != nil {
			return err
		}
	}

	// t.WinPoStProof ([]block.PoStProof) (slice)
	if len(t.WinPoStProof) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.WinPoStProof was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.WinPoStProof))); err != nil {
		return err
	}
	for _, v := range t.WinPoStProof {
		if err := v.MarshalCBOR(w); err != nil {
			return err
		}
	}

	// t.Parents (block.TipSetKey) (struct)
	if err := t.Parents.MarshalCBOR(w); err != nil {
		return err
	}

	// t.ParentWeight (big.Int) (struct)
	if err := t.ParentWeight.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Height (abi.ChainEpoch) (int64)
	if t.Height >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.Height)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.Height-1)); err != nil {
			return err
		}
	}

	// t.ParentStateRoot (cid.Cid) (struct)

	if err := cbg.WriteCidBuf(scratch, w, t.ParentStateRoot); err != nil {
		return xerrors.Errorf("failed to write cid field t.ParentStateRoot: %w", err)
	}

	// t.ParentMessageReceipts (cid.Cid) (struct)

	if err := cbg.WriteCidBuf(scratch, w, t.ParentMessageReceipts); err != nil {
		return xerrors.Errorf("failed to write cid field t.ParentMessageReceipts: %w", err)
	}

	// t.Messages (cid.Cid) (struct)

	if err := cbg.WriteCidBuf(scratch, w, t.Messages); err != nil {
		return xerrors.Errorf("failed to write cid field t.Messages: %w", err)
	}

	// t.BLSAggregate (crypto.Signature) (struct)
	if err := t.BLSAggregate.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Timestamp (uint64) (uint64)

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.Timestamp)); err != nil {
		return err
	}

	// t.BlockSig (crypto.Signature) (struct)
	if err := t.BlockSig.MarshalCBOR(w); err != nil {
		return err
	}

	// t.ForkSignaling (uint64) (uint64)

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.ForkSignaling)); err != nil {
		return err
	}

	// t.ParentBaseFee (big.Int) (struct)
	if err := t.ParentBaseFee.MarshalCBOR(w); err != nil {
		return err
	}
	return nil
}
