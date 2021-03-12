package core

import (
	"errors"
	"github.com/ahmetb/go-linq/v3"
	"github.com/filecoin-project/specs-actors/v3/actors/builtin"
	"github.com/ipfs-force-community/venus-wallet/errcode"
	"math"
	"reflect"
)

var (
	ErrCodeOverflow = errors.New("code over flow")
)

const (
	MethodSend        = builtin.MethodSend
	MethodConstructor = builtin.MethodConstructor
)

type MsgEnum = uint32

const (
	MEUnknown MsgEnum = 1 << iota
	MEChainMsg
	MEBlock
	MEDealProposal
	MEDeals
	MEDrawRandomParam
	MESignedVoucher
	MEStorageAsk
	MEAskResponse
	MENetWorkResponse
	MEProviderDealState
	MEClientDeal
)

var MsgEnumPool = []struct {
	Code int
	Name string
}{
	{Code: MsgEnumCode(MEUnknown), Name: "unknown"},
	{Code: MsgEnumCode(MEChainMsg), Name: "chainMsg"},
	{Code: MsgEnumCode(MEBlock), Name: "block"},
	{Code: MsgEnumCode(MEDealProposal), Name: "dealProposal"},
	{Code: MsgEnumCode(MEDeals), Name: "deals"},
	{Code: MsgEnumCode(MEDrawRandomParam), Name: "drawRandomParam"},
	{Code: MsgEnumCode(MESignedVoucher), Name: "signedVoucher"},
	{Code: MsgEnumCode(MEStorageAsk), Name: "storageAsk"},
	{Code: MsgEnumCode(MEAskResponse), Name: "askResponse"},
	{Code: MsgEnumCode(MENetWorkResponse), Name: "netWorkResponse"},
	{Code: MsgEnumCode(MEProviderDealState), Name: "providerDealState"},
	{Code: MsgEnumCode(MEClientDeal), Name: "clientDeal"},
}
var MaxMsgEnumCode = len(MsgEnumPool) - 1

func CheckMsgEnum(me MsgEnum) error {
	max := 1 << MaxMsgEnumCode
	if me < 0 || me > uint32(max) {
		return ErrCodeOverflow
	}
	return nil
}
func AggregateMsgEnumCode(codes []int) (MsgEnum, error) {
	if len(codes) == 0 {
		return 0, errcode.ErrNilReference
	}
	linq.From(codes).Distinct().ToSlice(codes)
	em := MsgEnum(0)
	for _, v := range codes {
		code, err := MsgEnumFromInt(v)
		if err != nil {
			return 0, err
		}
		em += 1 << code
	}
	return em, nil
}

func MsgEnumFromInt(code int) (MsgEnum, error) {
	if code < 0 || code > MaxMsgEnumCode {
		return 0, ErrCodeOverflow
	}
	return 1 << code, nil
}

func MsgEnumCode(me MsgEnum) int {
	code := math.Log2(float64(me))
	return int(code)
}
func ContainMsgType(multiME MsgEnum, mt MsgType) bool {
	me := convertToMsgEnum(mt)
	return multiME&me == me
}

func convertToMsgEnum(mt MsgType) MsgEnum {
	switch mt {
	case MTUnknown:
		return MEUnknown
	case MTChainMsg:
		return MEChainMsg
	case MTBlock:
		return MEBlock
	case MTDealProposal:
		return MEDealProposal
	case MTDeals:
		return MEDeals
	case MTDrawRandomParam:
		return MEDrawRandomParam
	case MTSignedVoucher:
		return MESignedVoucher
	case MTStorageAsk:
		return MEStorageAsk
	case MTAskResponse:
		return MEAskResponse
	case MTNetWorkResponse:
		return MENetWorkResponse
	case MTProviderDealState:
		return MEProviderDealState
	case MTClientDeal:
		return MEClientDeal
	default:
		return MEUnknown
	}
}

type MethodMeta struct {
	Name string

	Params reflect.Type
	Ret    reflect.Type
}
