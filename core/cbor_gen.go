package core

import (
	"bytes"
	"fmt"
	"github.com/filecoin-project/go-state-types/abi"
	proof2 "github.com/filecoin-project/specs-actors/v2/actors/runtime/proof"
	cbg "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"
	"io"
)

func (t *Block) UnmarshalCBOR(r io.Reader) error {
	*t = Block{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 16 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Miner (address.Address) (struct)

	{

		if err := t.Miner.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Miner: %w", err)
		}

	}
	// t.Ticket (newBlock.Ticket) (struct)

	{

		if err := t.Ticket.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Ticket: %w", err)
		}

	}
	// t.ElectionProof (newBlock.ElectionProof) (struct)

	{

		b, err := br.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := br.UnreadByte(); err != nil {
				return err
			}
			t.ElectionProof = new(ElectionProof)
			if err := t.ElectionProof.UnmarshalCBOR(br); err != nil {
				return xerrors.Errorf("unmarshaling t.ElectionProof pointer: %w", err)
			}
		}

	}
	// t.BeaconEntries ([]*newBlock.BeaconEntry) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.BeaconEntries: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.BeaconEntries = make([]*BeaconEntry, extra)
	}

	for i := 0; i < int(extra); i++ {

		var v BeaconEntry
		if err := v.UnmarshalCBOR(br); err != nil {
			return err
		}

		t.BeaconEntries[i] = &v
	}

	// t.WinPoStProof ([]newBlock.PoStProof) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.WinPoStProof: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.WinPoStProof = make([]proof2.PoStProof, extra)
	}

	for i := 0; i < int(extra); i++ {

		var v proof2.PoStProof
		if err := v.UnmarshalCBOR(br); err != nil {
			return err
		}

		t.WinPoStProof[i] = v
	}

	// t.Parents (newBlock.TipSetKey) (struct)

	{

		if err := t.Parents.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Parents: %w", err)
		}

	}
	// t.ParentWeight (big.Int) (struct)

	{

		if err := t.ParentWeight.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.ParentWeight: %w", err)
		}

	}
	// t.Height (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
		var extraI int64
		if err != nil {
			return err
		}
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative oveflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Height = abi.ChainEpoch(extraI)
	}
	// t.ParentStateRoot (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.ParentStateRoot: %w", err)
		}

		t.ParentStateRoot = c

	}
	// t.ParentMessageReceipts (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.ParentMessageReceipts: %w", err)
		}

		t.ParentMessageReceipts = c

	}
	// t.Messages (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Messages: %w", err)
		}

		t.Messages = c

	}
	// t.BLSAggregate (crypto.Signature) (struct)

	{

		b, err := br.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := br.UnreadByte(); err != nil {
				return err
			}
			t.BLSAggregate = new(Signature)
			if err := t.BLSAggregate.UnmarshalCBOR(br); err != nil {
				return xerrors.Errorf("unmarshaling t.BLSAggregate pointer: %w", err)
			}
		}

	}
	// t.Timestamp (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Timestamp = uint64(extra)

	}
	// t.BlockSig (crypto.Signature) (struct)

	{

		b, err := br.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := br.UnreadByte(); err != nil {
				return err
			}
			t.BlockSig = new(Signature)
			if err := t.BlockSig.UnmarshalCBOR(br); err != nil {
				return xerrors.Errorf("unmarshaling t.BlockSig pointer: %w", err)
			}
		}

	}
	// t.ForkSignaling (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.ForkSignaling = uint64(extra)

	}
	// t.ParentBaseFee (big.Int) (struct)

	{

		if err := t.ParentBaseFee.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.ParentBaseFee: %w", err)
		}

	}
	return nil
}

func (t *Ticket) UnmarshalCBOR(r io.Reader) error {
	*t = Ticket{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.VRFProof (newBlock.VRFPi) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.VRFProof: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.VRFProof = make([]uint8, extra)
	}

	if _, err := io.ReadFull(br, t.VRFProof[:]); err != nil {
		return err
	}
	return nil
}

func (t *ElectionProof) UnmarshalCBOR(r io.Reader) error {
	*t = ElectionProof{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.WinCount (int64) (int64)
	{
		maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
		var extraI int64
		if err != nil {
			return err
		}
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative oveflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.WinCount = int64(extraI)
	}
	// t.VRFProof (newBlock.VRFPi) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.VRFProof: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.VRFProof = make([]uint8, extra)
	}

	if _, err := io.ReadFull(br, t.VRFProof[:]); err != nil {
		return err
	}
	return nil
}
func (t *BeaconEntry) UnmarshalCBOR(r io.Reader) error {
	*t = BeaconEntry{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Round (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Round = uint64(extra)

	}
	// t.Data ([]uint8) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.Data: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.Data = make([]uint8, extra)
	}

	if _, err := io.ReadFull(br, t.Data[:]); err != nil {
		return err
	}
	return nil
}

func (tipsetKey *TipSetKey) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)
	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Parents: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		cids := make([]Cid, extra)
		for i := 0; i < int(extra); i++ {

			c, err := cbg.ReadCid(br)
			if err != nil {
				return xerrors.Errorf("reading cid field t.Parents failed: %v", err)
			}
			cids[i] = c
		}
		tipsetKey.value = string(encodeKey(cids))
	}
	return nil
}
func encodeKey(cids []Cid) []byte {
	buffer := new(bytes.Buffer)
	for _, c := range cids {
		// bytes.Buffer.Write() err is documented to be always nil.
		_, _ = buffer.Write(c.Bytes())
	}
	return buffer.Bytes()
}
