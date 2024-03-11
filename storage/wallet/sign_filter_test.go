package wallet

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/venus-wallet/config"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus/venus-shared/testutil"
	"github.com/filecoin-project/venus/venus-shared/types"
)

func TestSignFilter_CheckSignMsg(t *testing.T) {
	ctx := context.Background()
	addr, err := address.NewFromString("f3uyk4vweulsdbeqfnx7g4swk2zaa4p5xnmcuqvecyuwoggvlfagruxippti2v7sc2lzyop72pyrkr2ks2xc7q")
	assert.NoError(t, err)

	header := &types.Message{}
	testutil.Provide(t, header)
	toAddr, err := address.NewIDAddress(1001)
	assert.NoError(t, err)
	header.From = addr
	header.To = toAddr
	header.Method = 32

	signMsg := SignMsg{
		SignType: types.MTChainMsg,
		Data:     header,
	}

	t.Run("pass", func(t *testing.T) {
		filter := NewSignFilter(&config.SignFilter{Expr: "jq -e '.SignType==\"message\" and .Data.To[1:]==\"01001\" and (.Data.Method == 32 or (.Data.Method >= 5  and .Data.Method <= 11 ) or (.Data.Method >= 18  and .Data.Method <= 20 ) or (.Data.Method >= 24  and .Data.Method <= 29 ))'"})
		assert.NoError(t, filter.CheckSignMsg(ctx, signMsg))
	})
	t.Run("not pass", func(t *testing.T) {
		filter := NewSignFilter(&config.SignFilter{Expr: "jq -e '.SignType==\"message\" and .Data.To[1:]==\"01002\" and (.Data.Method == 32 or (.Data.Method >= 5  and .Data.Method <= 11 ) or (.Data.Method >= 18  and .Data.Method <= 20 ) or (.Data.Method >= 24  and .Data.Method <= 29 ))'"})
		assert.NotNil(t, filter.CheckSignMsg(ctx, signMsg))
	})
}
