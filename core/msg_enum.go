package core

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/rt"
	exported0 "github.com/filecoin-project/specs-actors/actors/builtin/exported"
	exported2 "github.com/filecoin-project/specs-actors/v2/actors/builtin/exported"
	"github.com/filecoin-project/specs-actors/v3/actors/builtin"
	exported3 "github.com/filecoin-project/specs-actors/v3/actors/builtin/exported"
	"reflect"
	"runtime"
	"strings"
)

const (
	MethodSend        = builtin.MethodSend
	MethodConstructor = builtin.MethodConstructor
)

type MsgEnum uint64

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

var MethodsMap = map[Cid]map[abi.MethodNum]MethodMeta{}

func init() {
	// TODO: combine with the runtime actor registry.
	var actors []rt.VMActor
	actors = append(actors, exported0.BuiltinActors()...)
	actors = append(actors, exported2.BuiltinActors()...)
	actors = append(actors, exported3.BuiltinActors()...)

	for _, actor := range actors {
		exports := actor.Exports()
		methods := make(map[abi.MethodNum]MethodMeta, len(exports))

		// Explicitly add send, it's special.
		methods[builtin.MethodSend] = MethodMeta{
			Name:   "Send",
			Params: reflect.TypeOf(new(abi.EmptyValue)),
			Ret:    reflect.TypeOf(new(abi.EmptyValue)),
		}

		// Iterate over exported methods. Some of these _may_ be nil and
		// must be skipped.
		for number, export := range exports {
			if export == nil {
				continue
			}

			ev := reflect.ValueOf(export)
			et := ev.Type()

			// Extract the method names using reflection. These
			// method names always match the field names in the
			// `builtin.Method*` structs (tested in the specs-actors
			// tests).
			fnName := runtime.FuncForPC(ev.Pointer()).Name()
			fnName = strings.TrimSuffix(fnName[strings.LastIndexByte(fnName, '.')+1:], "-fm")
			switch abi.MethodNum(number) {
			case MethodSend:
				panic("method 0 is reserved for Send")
			case MethodConstructor:
				if fnName != "Constructor" {
					panic("method 1 is reserved for Constructor")
				}
			}
			methods[abi.MethodNum(number)] = MethodMeta{
				Name:   fnName,
				Params: et.In(1),
				Ret:    et.Out(0),
			}
		}
		MethodsMap[actor.Code()] = methods
	}
}
