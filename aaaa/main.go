package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-wallet/api/remotecli"
	"github.com/filecoin-project/venus-wallet/api/remotecli/httpparse"
	"github.com/filecoin-project/venus-wallet/core"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	url := "http://127.0.0.1:4678/rpc/v0"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIl19.kJo0y8G-iTsNMY1aDNLLwwOFpCLuLmYMpnkss0uU48I_6c12e1af-7529-47b3-b0e6-1b0c42a79a37"
	headers := http.Header{}
	headers.Add(httpparse.ServiceToken, "Bearer "+string(token))
	client, closer, err := remotecli.NewFullNodeRPC(ctx, url, headers)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer closer()
	addrs, err := client.WalletList(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	to, err := address.NewFromString("t01561")
	if err != nil {
		log.Fatal(err)
		return
	}

	msg := core.Message{
		Version:    0,
		To:         to,
		From:       addrs[0],
		Nonce:      0,
		Value:      core.TokenAmount{},
		GasLimit:   0,
		GasFeeCap:  core.TokenAmount{},
		GasPremium: core.TokenAmount{},
		Method:     26,
		Params:     nil,
	}
	blk, _ := msg.ToStorageBlock()
	id := msg.Cid()
	xx, err := client.WalletSign(ctx, addrs[0], id.Bytes(), core.MsgMeta{
		Type:  core.MTChainMsg,
		Extra: blk.RawData(),
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(xx)
}
