package sqlite

import (
	"github.com/ipfs-force-community/venus-wallet/storage"
	"gotest.tools/assert"
	"testing"
)

func TestRouterStore_PutMsgTypeTemplate(t *testing.T) {
	mttName := "mockTest"
	mtt := &storage.MsgTypeTemplate{
		Name:      mttName,
		MetaTypes: 127,
	}
	err := mockRouterStore.PutMsgTypeTemplate(mtt)
	if err != nil {
		t.Fatalf("PutMsgTypeTemplate err:%s", err)
	}
	mttByName, err := mockRouterStore.GetMsgTypeTemplateByName(mttName)
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

}
