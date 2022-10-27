package core

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/ahmetb/go-linq/v3"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/venus-wallet/errcode"
	types "github.com/filecoin-project/venus/venus-shared/types/wallet"
	"github.com/filecoin-project/venus/venus-shared/utils"
	"github.com/ipfs/go-cid"
)

var (
	MethodNamesMap = make(map[string]struct{})
	MethodNameList []types.MethodName
)

type EmptyValue struct{}

func init() {
	loadMethodNames()
}

func loadMethodNames() {
	for _, methods := range utils.MethodsMap {
		for _, mm := range methods {
			MethodNamesMap[mm.Name] = struct{}{}
		}
	}

	MethodNameList = make([]types.MethodName, 0, len(MethodNamesMap))
	for k := range MethodNamesMap {
		MethodNameList = append(MethodNameList, k)
	}
	sort.Slice(MethodNameList, func(i, j int) bool {
		return MethodNameList[i] < MethodNameList[j]
	})
}

func ReloadMethodNames() {
	loadMethodNames()
}

func GetMethodName(actCode cid.Cid, method abi.MethodNum) (string, error) {
	m, found := utils.MethodsMap[actCode][method]
	if !found {
		return "", fmt.Errorf("unknown method %d for actor %s", method, actCode)
	}
	return m.Name, nil
}

func AggregateMethodNames(methods []types.MethodName) ([]types.MethodName, error) {
	if len(methods) == 0 {
		return nil, errcode.ErrNilReference
	}
	linq.From(methods).Distinct().ToSlice(&methods)
	var illegal []types.MethodName
	linq.From(methods).Except(linq.From(MethodNameList)).ToSlice(&illegal)
	buf := new(bytes.Buffer)
	if len(illegal) > 0 {
		for k, v := range illegal {
			if k > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(v)
		}
		return nil, fmt.Errorf("method name illegal: %s", buf.String())
	}
	return methods, nil
}
