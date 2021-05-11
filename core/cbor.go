package core

import (
	"bytes"
	"fmt"
	proof2 "github.com/filecoin-project/specs-actors/v2/actors/runtime/proof"
	"golang.org/x/xerrors"
	"io"

	"github.com/filecoin-project/go-state-types/abi"
	cbg "github.com/whyrusleeping/cbor-gen"
)

var lengthBufMessage = []byte{138}

func (msg *Message) MarshalCBOR(w io.Writer) error {
	if msg == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufMessage); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Version (uint64) (uint64)
	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, msg.Version); err != nil {
		return err
	}

	// t.To (address.Address) (struct)
	if err := msg.To.MarshalCBOR(w); err != nil {
		return err
	}

	// t.From (address.Address) (struct)
	if err := msg.From.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Nonce (uint64) (uint64)
	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, msg.Nonce); err != nil {
		return err
	}

	// t.Value (big.Int) (struct)
	if err := msg.Value.MarshalCBOR(w); err != nil {
		return err
	}

	// t.GasLimit (int64) (int64)
	if msg.GasLimit >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(msg.GasLimit)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-msg.GasLimit-1)); err != nil {
			return err
		}
	}

	// t.GasFeeCap (big.Int) (struct)
	if err := msg.GasFeeCap.MarshalCBOR(w); err != nil {
		return err
	}

	// t.GasPremium (big.Int) (struct)
	if err := msg.GasPremium.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Method (abi.MethodNum) (uint64)
	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(msg.Method)); err != nil {
		return err
	}

	// t.Params ([]uint8) (slice)
	if len(msg.Params) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Params was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(msg.Params))); err != nil {
		return err
	}

	if _, err := w.Write(msg.Params[:]); err != nil {
		return err
	}
	return nil
}

func (msg *Message) UnmarshalCBOR(r io.Reader) error {
	*msg = Message{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 10 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Version (uint64) (uint64)
	{
		maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		msg.Version = extra
	}
	// t.To (address.Address) (struct)
	{
		if err := msg.To.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.To: %w", err)
		}
	}
	// t.From (address.Address) (struct)

	{
		if err := msg.From.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.From: %w", err)
		}
	}
	// t.Nonce (uint64) (uint64)
	{
		maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		msg.Nonce = extra
	}
	// t.Value (big.Int) (struct)
	{
		if err := msg.Value.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Value: %w", err)
		}
	}
	// t.GasLimit (int64) (int64)
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
		msg.GasLimit = extraI
	}
	// t.GasFeeCap (big.Int) (struct)
	{
		if err := msg.GasFeeCap.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.GasFeeCap: %w", err)
		}
	}
	// t.GasPremium (big.Int) (struct)
	{
		if err := msg.GasPremium.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.GasPremium: %w", err)
		}
	}
	// t.Method (abi.MethodNum) (uint64)
	{
		maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		msg.Method = abi.MethodNum(extra)
	}
	// t.Params ([]uint8) (slice)
	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.Params: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		msg.Params = make([]uint8, extra)
	}

	if _, err := io.ReadFull(br, msg.Params[:]); err != nil {
		return err
	}
	return nil
}

var lengthBufSignedMessage = []byte{130}

func (t *SignedMessage) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufSignedMessage); err != nil {
		return err
	}

	// t.Message (types.Message) (struct)
	if err := t.Message.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Signature (crypto.Signature) (struct)
	if err := t.Signature.MarshalCBOR(w); err != nil {
		return err
	}
	return nil
}

func (t *SignedMessage) UnmarshalCBOR(r io.Reader) error {
	*t = SignedMessage{}

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

	// t.Message (types.Message) (struct)

	{

		if err := t.Message.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Message: %w", err)
		}

	}
	// t.Signature (crypto.Signature) (struct)

	{

		if err := t.Signature.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Signature: %w", err)
		}

	}
	return nil
}

func (bk *Block) UnmarshalCBOR(r io.Reader) error {
	*bk = Block{}

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
		if err := bk.Miner.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Miner: %w", err)
		}
	}
	// t.Ticket (newBlock.Ticket) (struct)
	{
		if err := bk.Ticket.UnmarshalCBOR(br); err != nil {
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
			bk.ElectionProof = new(ElectionProof)
			if err := bk.ElectionProof.UnmarshalCBOR(br); err != nil {
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
		bk.BeaconEntries = make([]*BeaconEntry, extra)
	}

	for i := 0; i < int(extra); i++ {
		var v BeaconEntry
		if err := v.UnmarshalCBOR(br); err != nil {
			return err
		}
		bk.BeaconEntries[i] = &v
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
		bk.WinPoStProof = make([]proof2.PoStProof, extra)
	}

	for i := 0; i < int(extra); i++ {
		var v proof2.PoStProof
		if err := v.UnmarshalCBOR(br); err != nil {
			return err
		}
		bk.WinPoStProof[i] = v
	}

	// t.Parents (newBlock.TipSetKey) (struct)
	{
		if err := bk.Parents.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Parents: %w", err)
		}
	}
	// t.ParentWeight (big.Int) (struct)
	{
		if err := bk.ParentWeight.UnmarshalCBOR(br); err != nil {
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

		bk.Height = abi.ChainEpoch(extraI)
	}
	// t.ParentStateRoot (cid.Cid) (struct)
	{
		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.ParentStateRoot: %w", err)
		}

		bk.ParentStateRoot = c
	}
	// t.ParentMessageReceipts (cid.Cid) (struct)
	{
		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.ParentMessageReceipts: %w", err)
		}

		bk.ParentMessageReceipts = c
	}
	// t.Messages (cid.Cid) (struct)
	{
		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Messages: %w", err)
		}

		bk.Messages = c
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
			bk.BLSAggregate = new(Signature)
			if err := bk.BLSAggregate.UnmarshalCBOR(br); err != nil {
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
		bk.Timestamp = uint64(extra)
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
			bk.BlockSig = new(Signature)
			if err := bk.BlockSig.UnmarshalCBOR(br); err != nil {
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
		bk.ForkSignaling = uint64(extra)
	}
	// t.ParentBaseFee (big.Int) (struct)
	{
		if err := bk.ParentBaseFee.UnmarshalCBOR(br); err != nil {
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
