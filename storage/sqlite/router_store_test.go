package sqlite

import (
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/msgrouter"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"gotest.tools/assert"
	"math/rand"
	"testing"
)

func TestRouterStore_PutMsgTypeTemplate(t *testing.T) {
	mockName := "mockTest"
	// MsgTypeTemplate test
	mtt := &storage.MsgTypeTemplate{
		Name:      mockName,
		MetaTypes: 127,
	}
	err := mockRouterStore.PutMsgTypeTemplate(mtt)
	if err != nil {
		t.Fatalf("PutMsgTypeTemplate err:%s", err)
	}
	mttByName, err := mockRouterStore.GetMsgTypeTemplatesByName(mockName)
	if err != nil {
		t.Fatalf("GetMsgTypeTemplateByName err:%s", err)
	}
	assert.Equal(t, len(mttByName), 1)
	mttById, err := mockRouterStore.GetMsgTypeTemplate(mttByName[0].MTTId)
	if err != nil {
		t.Fatalf("GetMsgTypeTemplate err:%s", err)
	}
	assert.DeepEqual(t, mttById, mttByName[0])

	mttArr, err := mockRouterStore.ListMsgTypeTemplates(0, 10)
	if err != nil {
		t.Fatalf("ListMethodTemplates err:%s", err)
	}
	assert.Equal(t, len(mttArr), 1)
	assert.DeepEqual(t, mttById, mttArr[0])

	err = mockRouterStore.DeleteMsgTypeTemplate(mttById.MTTId)
	if err != nil {
		t.Fatalf("DeleteMsgTypeTemplate err:%s", err)
	}

	// MethodTemplate test
	mtCount := 10
	source := msgrouter.MethodNameList
	//random name
	for i := 0; i < mtCount; i++ {
		var methodNames []string
		for j := 0; j < 10; j++ {
			idx := rand.Intn(len(source))
			methodNames = append(methodNames, source[idx])
		}
		methodNames = Unique(methodNames)
		mt := &storage.MethodTemplate{
			Name:    mockName,
			Methods: methodNames,
		}
		err = mockRouterStore.PutMethodTemplate(mt)
		if err != nil {
			t.Fatalf("The serial number：%d PutMethodTemplate error:%s", i, err)
		}
	}
	mtArrByName, err := mockRouterStore.GetMethodTemplatesByName(mockName)
	if err != nil {
		t.Fatalf("GetMsgTypeTemplateByName err:%s", err)
	}
	assert.Equal(t, len(mtArrByName), mtCount)
	mtById, err := mockRouterStore.GetMethodTemplate(mtArrByName[0].MTId)
	if err != nil {
		t.Fatalf("GetMsgTypeTemplate err:%s", err)
	}
	assert.DeepEqual(t, mtById, mtArrByName[0])
	mtArr, err := mockRouterStore.ListMethodTemplates(0, mtCount)
	if err != nil {
		t.Fatalf("ListMethodTemplates err:%s", err)
	}
	assert.DeepEqual(t, mtArr, mtArrByName)
	for k, v := range mtArr {
		err = mockRouterStore.DeleteMethodTemplate(v.MTId)
		if err != nil {
			t.Fatalf("The serial number：%d DeleteMethodTemplate err:%s", k, err)
		}
	}

	// keyBind test
	mockAddress := "f3vnzhpj6xftqubkp2klnce37jpyndwurj4eyoqs7sx3weegth5joxmphtos4ni6tuxmidc2nj55ygag33qesq"
	for i := 0; i < mtCount; i++ {
		var methodNames []string
		for j := 0; j < 10; j++ {
			idx := rand.Intn(len(source))
			methodNames = append(methodNames, source[idx])
		}
		methodNames = Unique(methodNames)
		metaTypes := core.MsgEnum(rand.Intn(30))
		kb := &storage.KeyBind{
			Name:      mockName,
			MetaTypes: metaTypes,
			Address:   mockAddress,
			Methods:   methodNames,
		}
		err = mockRouterStore.PutKeyBind(kb)
		if err != nil {
			t.Fatalf("The serial number：%d PutKeyBind error:%s", i, err)
		}
	}
	kbArrByAddress, err := mockRouterStore.GetKeyBinds(mockAddress)
	if err != nil {
		t.Fatalf("GetKeyBinds err:%s", err)
	}
	assert.Equal(t, len(kbArrByAddress), mtCount)
	kbById, err := mockRouterStore.GetKeyBindById(kbArrByAddress[0].BindId)
	if err != nil {
		t.Fatalf("GetKeyBindById err:%s", err)
	}
	assert.DeepEqual(t, kbById, kbArrByAddress[0])
	kbArr, err := mockRouterStore.ListKeyBinds(0, mtCount)
	if err != nil {
		t.Fatalf("ListKeyBinds err:%s", err)
	}
	assert.DeepEqual(t, kbArr, kbArrByAddress)
	kbArrByName, err := mockRouterStore.GetKeyBindsByName(mockName)
	if err != nil {
		t.Fatalf("GetKeyBindsByName err:%s", err)
	}
	assert.DeepEqual(t, kbArrByName, kbArrByAddress)
	for k, v := range kbArrByName {
		if k*2 >= mtCount {
			break
		}
		err = mockRouterStore.DeleteKeyBind(v.BindId)
		if err != nil {
			t.Fatalf("The serial number：%d DeleteKeyBind err:%s", k, err)
		}
	}
	delCount, err := mockRouterStore.DeleteKeyBindsByAddress(mockAddress)
	if err != nil {
		t.Fatalf("DeleteKeyBindsByAddress err:%s", err)
	}
	assert.Equal(t, delCount, int64(mtCount/2))

	// Group Test

}

func Unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := make([]string, 0)
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func generateGroupPreSource(count int) {

}
