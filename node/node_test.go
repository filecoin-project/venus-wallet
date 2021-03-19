package node

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/venus-wallet/config"
	"os"
	"testing"
)

func TestLotusNode_GetActorCode(t *testing.T) {
	if os.Getenv("CI") == "test" {
		t.Skip()
	}
	cli, err := NewNodeClient(&config.StrategyConfig{
		Level:   3,
		NodeURL: "/ip4/127.0.0.1/tcp/1234/http",
	})
	if err != nil {
		t.Fatal(err)
	}
	addr, _ := address.NewFromString("t3spqxxmzgmz2t6flflnw23c2shdk45thqxt5dxyoacwuawe6bmt37xicxdeqw2fwxqs2mtdnfweqfjcome7ka")
	actor, err := cli.StateGetActor(context.Background(), addr, TipSetKey{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(actor.Code)
}
