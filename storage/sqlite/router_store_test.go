package sqlite

import (
	"github.com/google/uuid"
	"github.com/ipfs-force-community/venus-wallet/core"
	"github.com/ipfs-force-community/venus-wallet/msgrouter"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"gotest.tools/assert"
	"math/rand"
	"strconv"
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
	mttByName, err := mockRouterStore.GetMsgTypeTemplateByName(mockName)
	if err != nil {
		t.Fatalf("GetMsgTypeTemplateByName err:%s", err)
	}
	mttById, err := mockRouterStore.GetMsgTypeTemplate(mttByName.MTTId)
	if err != nil {
		t.Fatalf("GetMsgTypeTemplate err:%s", err)
	}
	assert.DeepEqual(t, mttById, mttByName)

	mttArr, err := mockRouterStore.ListMsgTypeTemplates(0, 10)
	if err != nil {
		t.Fatalf("ListMethodTemplates err:%s", err)
	}
	assert.Equal(t, len(mttArr), 1)
	assert.DeepEqual(t, mttById, mttArr[0])

	// MethodTemplate test
	mtCount := 10
	source := msgrouter.MethodNameList
	mtArrByName := make([]*storage.MethodTemplate, 0)
	//random name
	for i := 0; i < mtCount; i++ {
		var methodNames []string
		for j := 0; j < 10; j++ {
			idx := rand.Intn(len(source))
			methodNames = append(methodNames, source[idx])
		}
		methodNames = Unique(methodNames)
		mockNameI := mockName + strconv.Itoa(i)
		mt := &storage.MethodTemplate{
			Name:    mockNameI,
			Methods: methodNames,
		}
		err = mockRouterStore.PutMethodTemplate(mt)
		if err != nil {
			t.Fatalf("The serial number：%d PutMethodTemplate error:%s", i, err)
		}
		mtTmp, err := mockRouterStore.GetMethodTemplateByName(mockNameI)
		if err != nil {
			t.Fatalf("GetMsgTypeTemplateByName err:%s", err)
		}
		mtArrByName = append(mtArrByName, mtTmp)
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

	// keyBind test
	kbArrByName := make([]*storage.KeyBind, 0)
	mockAddress := "f3vnzhpj6xftqubkp2klnce37jpyndwurj4eyoqs7sx3weegth5joxmphtos4ni6tuxmidc2nj55ygag33qesq"
	for i := 0; i < mtCount; i++ {
		var methodNames []string
		for j := 0; j < 10; j++ {
			idx := rand.Intn(len(source))
			methodNames = append(methodNames, source[idx])
		}
		methodNames = Unique(methodNames)
		mockNameI := mockName + strconv.Itoa(i)
		metaTypes := core.MsgEnum(rand.Intn(30))
		kb := &storage.KeyBind{
			Name:      mockNameI,
			MetaTypes: metaTypes,
			Address:   mockAddress,
			Methods:   methodNames,
		}
		err = mockRouterStore.PutKeyBind(kb)
		if err != nil {
			t.Fatalf("The serial number：%d PutKeyBind error:%s", i, err)
		}
		kbTmp, err := mockRouterStore.GetKeyBindByName(mockNameI)
		if err != nil {
			t.Fatalf("GetKeyBindsByName err:%s", err)
		}
		kbArrByName = append(kbArrByName, kbTmp)
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
	assert.DeepEqual(t, kbArrByName, kbArrByAddress)

	// Group Test
	mockGroupName := "mockGroup"
	err = mockRouterStore.PutGroup(mockGroupName, []uint{kbById.BindId})
	if err != nil {
		t.Fatalf("PutGroup err:%s", err)
	}
	groupByName, err := mockRouterStore.GetGroupByName(mockGroupName)
	if err != nil {
		t.Fatalf("GetGroupByName err:%s", err)
	}
	groupById, err := mockRouterStore.GetGroup(groupByName.GroupId)
	if err != nil {
		t.Fatalf("GetGroup err:%s", err)
	}
	assert.DeepEqual(t, groupByName, groupById)

	// groupAuth
	token := uuid.New().String()
	err = mockRouterStore.PutGroupAuth(token, groupById.GroupId)
	if err != nil {
		t.Fatalf("PutGroupAuth err:%s", err)
	}
	gAuth, err := mockRouterStore.GetGroupAuth(token)
	if err != nil {
		t.Fatalf("GetGroupAuth err:%s", err)
	}
	assert.DeepEqual(t, gAuth.KeyBinds, groupById.KeyBinds)

	// release all
	// MsgTypeTemplate DEL
	err = mockRouterStore.DeleteMsgTypeTemplate(mttById.MTTId)
	if err != nil {
		t.Fatalf("DeleteMsgTypeTemplate err:%s", err)
	}

	for k, v := range mtArr {
		err = mockRouterStore.DeleteMethodTemplate(v.MTId)
		if err != nil {
			t.Fatalf("The serial number：%d DeleteMethodTemplate err:%s", k, err)
		}
	}

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

	err = mockRouterStore.DeleteGroup(groupById.GroupId)
	if err != nil {
		t.Fatalf("DeleteKeyBindsByAddress err:%s", err)
	}
	err = mockRouterStore.DeleteGroupAuth(token)
	if err != nil {
		t.Fatalf("DeleteGroupAuth err:%s", err)
	}
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
