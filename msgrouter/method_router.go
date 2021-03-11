package msgrouter

import (
	"fmt"
	"github.com/filecoin-project/go-state-types/rt"
	exported0 "github.com/filecoin-project/specs-actors/actors/builtin/exported"
	exported2 "github.com/filecoin-project/specs-actors/v2/actors/builtin/exported"
	"github.com/filecoin-project/specs-actors/v3/actors/builtin"
	exported3 "github.com/filecoin-project/specs-actors/v3/actors/builtin/exported"
	"github.com/ipfs-force-community/venus-wallet/core"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

var MethodsMap = map[core.Cid]map[core.MethodNum]core.MethodMeta{}

var MethodNamesMap = make(map[string]struct{})
var MethodNameList []MethodName

type MethodName = string
type EmptyValue struct{}

func init() {
	// TODO: combine with the runtime actor registry.
	var actors []rt.VMActor
	actors = append(actors, exported0.BuiltinActors()...)
	actors = append(actors, exported2.BuiltinActors()...)
	actors = append(actors, exported3.BuiltinActors()...)

	for _, actor := range actors {
		exports := actor.Exports()
		methods := make(map[core.MethodNum]core.MethodMeta, len(exports))

		// Explicitly add send, it's special.
		methods[builtin.MethodSend] = core.MethodMeta{
			Name:   "Send",
			Params: reflect.TypeOf(new(EmptyValue)),
			Ret:    reflect.TypeOf(new(EmptyValue)),
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
			switch core.MethodNum(number) {
			case core.MethodSend:
				panic("method 0 is reserved for Send")
			case core.MethodConstructor:
				if fnName != "Constructor" {
					panic("method 1 is reserved for Constructor")
				}
			}
			methods[core.MethodNum(number)] = core.MethodMeta{
				Name:   fnName,
				Params: et.In(1),
				Ret:    et.Out(0),
			}
			MethodNamesMap[fnName] = struct{}{}
		}
		MethodsMap[actor.Code()] = methods
	}
	MethodNameList = make([]MethodName, 0, len(MethodNamesMap))
	for k, _ := range MethodNamesMap {
		MethodNameList = append(MethodNameList, k)
	}
	sort.Slice(MethodNameList, func(i, j int) bool {
		return MethodNameList[i] < MethodNameList[j]
	})
}

func GetMethodName(actCode core.Cid, method core.MethodNum) (string, error) {
	m, found := MethodsMap[actCode][method]
	if !found {
		return core.StringEmpty, fmt.Errorf("unknown method %d for actor %s", method, actCode)
	}
	return m.Name, nil
}
