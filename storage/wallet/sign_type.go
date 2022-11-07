package wallet

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/ipfs-force-community/venus-gateway/walletevent"

	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket/migrations"
	"github.com/filecoin-project/go-fil-markets/storagemarket/network"
	"github.com/filecoin-project/go-state-types/cbor"
	"github.com/filecoin-project/specs-actors/v2/actors/builtin/market"
	"github.com/filecoin-project/specs-actors/v2/actors/builtin/paych"
	"github.com/filecoin-project/venus/venus-shared/types"
	types2 "github.com/filecoin-project/venus/venus-shared/types/wallet"
)

// Types Abstract data types to be signed
type Types struct {
	Type      reflect.Type
	signBytes FGetSignBytes
	parseObj  FParseObj
}

type (
	FGetSignBytes func(in interface{}) ([]byte, error)
	FParseObj     func([]byte, types.MsgMeta) (interface{}, error)
)

func RegisterSupportedMsgTypes(msgType types.MsgType, p reflect.Type,
	fGetSignBytes FGetSignBytes, fParseObj FParseObj,
) (replaced bool) {
	_, replaced = SupportedMsgTypes[msgType]
	SupportedMsgTypes[msgType] = &Types{p, fGetSignBytes, fParseObj}
	return replaced
}

// SupportedMsgTypes signature type factory
var SupportedMsgTypes = map[types.MsgType]*Types{
	types.MTDealProposal: {reflect.TypeOf(market.DealProposal{}), func(i interface{}) ([]byte, error) {
		return cborutil.Dump(i)
	}, nil},
	types.MTClientDeal: {reflect.TypeOf(market.ClientDealProposal{}), func(in interface{}) ([]byte, error) {
		ni, err := cborutil.AsIpld(in)
		if err != nil {
			return nil, err
		}
		return ni.Cid().Bytes(), nil
	}, nil},
	types.MTDrawRandomParam: {reflect.TypeOf(types2.DrawRandomParams{}), func(in interface{}) ([]byte, error) {
		param := in.(*types2.DrawRandomParams)
		return param.SignBytes()
	}, nil},
	types.MTSignedVoucher: {reflect.TypeOf(paych.SignedVoucher{}), func(in interface{}) ([]byte, error) {
		return (in.(*paych.SignedVoucher)).SigningBytes()
	}, nil},
	types.MTStorageAsk: {reflect.TypeOf(storagemarket.StorageAsk{}), func(in interface{}) ([]byte, error) {
		return cborutil.Dump(in)
	}, nil},
	types.MTAskResponse: {Type: reflect.TypeOf(network.AskResponse{}), signBytes: func(in interface{}) ([]byte, error) {
		newAsk := in.(*network.AskResponse).Ask.Ask
		oldAsk := &migrations.StorageAsk0{
			Price: newAsk.Price, VerifiedPrice: newAsk.VerifiedPrice, MinPieceSize: newAsk.MinPieceSize,
			MaxPieceSize: newAsk.MaxPieceSize, Miner: newAsk.Miner, Timestamp: newAsk.Timestamp, Expiry: newAsk.Expiry, SeqNo: newAsk.SeqNo,
		}
		return cborutil.Dump(oldAsk)
	}},
	types.MTNetWorkResponse: {reflect.TypeOf(network.Response{}), func(in interface{}) ([]byte, error) {
		return cborutil.Dump(in)
	}, nil},

	types.MTBlock: {reflect.TypeOf(types.BlockHeader{}), func(in interface{}) ([]byte, error) {
		return in.(*types.BlockHeader).SignatureData()
	}, nil},
	types.MTChainMsg: {reflect.TypeOf(types.Message{}), func(in interface{}) ([]byte, error) {
		msg := in.(*types.Message)
		return msg.Cid().Bytes(), nil
	}, nil},
	types.MTProviderDealState: {
		reflect.TypeOf(storagemarket.ProviderDealState{}), func(in interface{}) ([]byte, error) {
			return cborutil.Dump(in)
		}, nil,
	},
	// chain/gen/gen.go:659,
	// in method 'ComputeVRF' sign bytes with MsgType='MTUnknown'
	// so, must deal 'MTUnknown' MsgType, and this may case safe problem
	types.MTUnknown: {reflect.TypeOf([]byte{}), func(in interface{}) ([]byte, error) {
		msg, isok := in.([]byte)
		if !isok {
			return nil, fmt.Errorf("MTUnkown must be []byte")
		}
		return msg, nil
	}, func(in []byte, meta types.MsgMeta) (interface{}, error) {
		if meta.Type == types.MTUnknown {
			return in, nil
		}
		return nil, fmt.Errorf("un-expected MsgType:%s", meta.Type)
	}},
	// the data to sign is divide into 2 parts:
	// first  part: is from venus-gateway, which here should be `meta.Extra`
	// second part: is from venus-wallet, which here is `wallet_event.RandomBytes`
	types.MTVerifyAddress: {
		Type: reflect.TypeOf([]byte{}),
		signBytes: func(in interface{}) ([]byte, error) {
			return in.([]byte), nil
		},
		parseObj: func(in []byte, meta types.MsgMeta) (interface{}, error) {
			expected := walletevent.GetSignData(meta.Extra, walletevent.RandomBytes)
			if !bytes.Equal(in, expected) {
				return nil, fmt.Errorf("sign data not match, actual %v, expected %v", in, expected)
			}
			return in, nil
		},
	},
}

// GetSignBytes Matches the type and returns the data that needs to be signed
func GetSignBytes(toSign []byte, meta types.MsgMeta) (interface{}, []byte, error) {
	t := SupportedMsgTypes[meta.Type]
	if t == nil {
		return nil, nil, fmt.Errorf("unsupported msgtype:%s", meta.Type)
	}
	var in interface{}
	var err error
	if t.parseObj != nil {
		if in, err = t.parseObj(toSign, meta); err != nil {
			return nil, nil, fmt.Errorf("parseObj failed:%w", err)
		}
	} else { // treat as cbor unmarshal-able object by default
		in = reflect.New(t.Type).Interface()
		unmarshaler, isok := in.(cbor.Unmarshaler)
		if !isok {
			return nil, nil, fmt.Errorf("type:%s is is not an 'unmarhsaler'", t.Type.Name())
		}
		if err := unmarshaler.UnmarshalCBOR(bytes.NewReader(toSign)); err != nil {
			return nil, nil, fmt.Errorf("cborunmarshal to %s failed:%w", t.Type.Name(), err)
		}
	}
	var data []byte
	data, err = t.signBytes(in)
	return in, data, err
}
