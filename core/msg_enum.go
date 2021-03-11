package core

import (
	"github.com/filecoin-project/specs-actors/v3/actors/builtin"
	"reflect"
)

const (
	MethodSend        = builtin.MethodSend
	MethodConstructor = builtin.MethodConstructor
)

type MsgEnum = uint32

const (
	MEUnknown           MsgEnum = 1 << 0
	MEChainMsg          MsgEnum = 1 << 1
	MEBlock             MsgEnum = 1 << 2
	MEDealProposal      MsgEnum = 1 << 3
	MEDeals             MsgEnum = 1 << 4
	MEDrawRandomParam   MsgEnum = 1 << 5
	MESignedVoucher     MsgEnum = 1 << 6
	MEStorageAsk        MsgEnum = 1 << 7
	MEAskResponse       MsgEnum = 1 << 8
	MENetWorkResponse   MsgEnum = 1 << 9
	MEProviderDealState MsgEnum = 1 << 10
	MEClientDeal        MsgEnum = 1 << 11
)

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
