package core

import (
	"gotest.tools/assert"
	"testing"
)

func TestContainMsgType(t *testing.T) {
	multiME := MEChainMsg + MEStorageAsk + MEProviderDealState
	assert.Equal(t, ContainMsgType(multiME, MTChainMsg), true)
	assert.Equal(t, ContainMsgType(multiME, MTStorageAsk), true)
	assert.Equal(t, ContainMsgType(multiME, MTProviderDealState), true)

	assert.Equal(t, ContainMsgType(multiME, MTUnknown), false)
	assert.Equal(t, ContainMsgType(multiME, MTBlock), false)
	assert.Equal(t, ContainMsgType(multiME, MTDealProposal), false)
	assert.Equal(t, ContainMsgType(multiME, MTDrawRandomParam), false)
	assert.Equal(t, ContainMsgType(multiME, MTSignedVoucher), false)
	assert.Equal(t, ContainMsgType(multiME, MTAskResponse), false)
	assert.Equal(t, ContainMsgType(multiME, MTNetWorkResponse), false)
	assert.Equal(t, ContainMsgType(multiME, MTClientDeal), false)
}

func TestFindCode(t *testing.T) {
	ids := FindCode(38)
	assert.DeepEqual(t, []int{1, 2, 5}, ids)

	ids2 := FindCode(8392)
	assert.DeepEqual(t, []int{3, 6, 7, 13}, ids2)
}

func TestAggregateMsgEnumCode(t *testing.T) {
	me, err := AggregateMsgEnumCode([]int{1, 2, 3, 4, 5, 6, 7})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, me, uint32(254))
}
