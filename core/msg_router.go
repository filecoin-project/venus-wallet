package core

import (
	"bytes"
	"fmt"
	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket/migrations"
	"github.com/filecoin-project/go-fil-markets/storagemarket/network"
	"github.com/filecoin-project/go-state-types/cbor"
	"github.com/filecoin-project/specs-actors/v2/actors/builtin/market"
	"github.com/filecoin-project/specs-actors/v2/actors/builtin/paych"
	"golang.org/x/xerrors"
	"reflect"
)

// Abstract data types to be signed
type Types struct {
	Type      reflect.Type
	signBytes func(i interface{}) ([]byte, error)
	parseObj  func([]byte, MsgMeta) (interface{}, error)
}

// signature type factory
var SupportedMsgTypes = map[MsgType]*Types{
	MTDealProposal: {reflect.TypeOf(market.DealProposal{}), func(i interface{}) ([]byte, error) {
		return cborutil.Dump(i)
	}, nil},
	MTDrawRandomParam: {reflect.TypeOf(DrawRandomParams{}), func(in interface{}) ([]byte, error) {
		param := in.(*DrawRandomParams)
		return param.SignBytes()
	}, nil},
	MTSignedVoucher: {reflect.TypeOf(paych.SignedVoucher{}), func(in interface{}) ([]byte, error) {
		return (in.(*paych.SignedVoucher)).SigningBytes()
	}, nil},
	MTStorageAsk: {reflect.TypeOf(storagemarket.StorageAsk{}), func(in interface{}) ([]byte, error) {
		return cborutil.Dump(in)
	}, nil},
	MTAskResponse: {Type: reflect.TypeOf(network.AskResponse{}), signBytes: func(in interface{}) ([]byte, error) {
		newAsk := in.(*network.AskResponse).Ask.Ask
		oldAsk := &migrations.StorageAsk0{Price: newAsk.Price, VerifiedPrice: newAsk.VerifiedPrice, MinPieceSize: newAsk.MinPieceSize,
			MaxPieceSize: newAsk.MaxPieceSize, Miner: newAsk.Miner, Timestamp: newAsk.Timestamp, Expiry: newAsk.Expiry, SeqNo: newAsk.SeqNo}
		return cborutil.Dump(oldAsk)
	}},
	MTNetWorkResponse: {reflect.TypeOf(network.Response{}), func(in interface{}) ([]byte, error) {
		return cborutil.Dump(in)
	}, nil},
	MTClientDeal: {reflect.TypeOf(market.ClientDealProposal{}), func(in interface{}) ([]byte, error) {
		ni, err := cborutil.AsIpld(in)
		if err != nil {
			return nil, err
		}
		return ni.Cid().Bytes(), nil
	}, nil},
	MTBlock: {reflect.TypeOf(Block{}), func(in interface{}) ([]byte, error) {
		return in.(*Block).SignatureData(), nil
	}, nil},
	MTChainMsg: {reflect.TypeOf(Message{}), func(in interface{}) ([]byte, error) {
		msg := in.(*Message)
		return msg.Cid().Bytes(), nil
	}, nil},
	// chain/gen/gen.go:659,
	// in method 'ComputVRF' sign bytes with MsgType='MTUnkown'
	// so, must deal 'MTUnkown' MsgType, and this may case safe problem
	MTUnknown: {reflect.TypeOf([]byte{}), func(in interface{}) ([]byte, error) {
		msg, isok := in.([]byte)
		if !isok {
			return nil, fmt.Errorf("MTUnkown must be []byte")
		}
		return msg, nil
	}, func(in []byte, meta MsgMeta) (interface{}, error) {
		if meta.Type == MTUnknown {
			return in, nil
		}
		return nil, fmt.Errorf("un-expected MsgType:%s", meta.Type)
	}},
}

// Matches the type and returns the data that needs to be signed
func GetSignBytes(toSign []byte, meta MsgMeta) (interface{}, []byte, error) {
	t := SupportedMsgTypes[meta.Type]
	if t == nil {
		return nil, nil, fmt.Errorf("unsupported msgtype:%s", meta.Type)
	}
	var in interface{}
	var err error
	if t.parseObj != nil {
		if in, err = t.parseObj(toSign, meta); err != nil {
			return nil, nil, xerrors.Errorf("parseObj failed:%w", err)
		}
	} else { // treat as cbor unmarshal-able object by default
		in = reflect.New(t.Type).Interface()
		unmarshaler, isok := in.(cbor.Unmarshaler)
		if !isok {
			return nil, nil, fmt.Errorf("type:%s is is not an 'unmarhsaler'", t.Type.Name())
		}
		if err := unmarshaler.UnmarshalCBOR(bytes.NewReader(toSign)); err != nil {
			return nil, nil, xerrors.Errorf("cborunmarshal to %s failed:%w", t.Type.Name(), err)
		}
	}
	var data []byte
	data, err = t.signBytes(in)
	return in, data, err
}
