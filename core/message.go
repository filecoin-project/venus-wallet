package core

import (
	"bytes"
	"fmt"
	"github.com/filecoin-project/go-state-types/abi"
	block "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
)

const MessageVersion = 0

// filecoin message
type Message struct {
	Version uint64

	To   Address
	From Address

	Nonce uint64

	Value TokenAmount

	GasLimit   int64
	GasFeeCap  TokenAmount
	GasPremium TokenAmount

	Method MethodNum
	Params []byte
}

func (m *Message) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := m.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Message) Cid() cid.Cid {
	b, err := m.ToStorageBlock()
	if err != nil {
		panic(fmt.Sprintf("failed to marshal message: %s", err)) // I think this is maybe sketchy, what happens if we try to serialize a message with an undefined address in it?
	}
	return b.Cid()
}

func (m *Message) ToStorageBlock() (block.Block, error) {
	data, err := m.Serialize()
	if err != nil {
		return nil, err
	}
	c, err := abi.CidBuilder.Sum(data)
	if err != nil {
		return nil, err
	}
	return block.NewBlockWithCid(data, c)
}

func DecodeMessage(b []byte) (*Message, error) {
	var msg Message
	if err := msg.UnmarshalCBOR(bytes.NewReader(b)); err != nil {
		return nil, err
	}

	if msg.Version != MessageVersion {
		return nil, fmt.Errorf("decoded message had incorrect version (%d)", msg.Version)
	}

	return &msg, nil
}

type SignedMessage struct {
	Message   Message
	Signature Signature
}
