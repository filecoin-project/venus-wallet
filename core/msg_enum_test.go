package core

import (
	"gotest.tools/assert"
	"testing"
)

func TestContainMsgType(t *testing.T) {
	multiME := MEChainMsg + MEDeals + MEStorageAsk + MEProviderDealState
	assert.Equal(t, ContainMsgType(multiME, MTChainMsg), true)
	assert.Equal(t, ContainMsgType(multiME, MTDeals), true)
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


