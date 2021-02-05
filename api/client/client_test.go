package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/filecoin-project/go-address"
	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/filecoin-project/go-fil-markets/shared_testutil"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket/network"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/node/impl/force/db_proc"
	api2 "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/ipfs/go-cid"
	"net/http"
	"os"
	"testing"
)

var walletAPI api.FullNode
var closer jsonrpc.ClientCloser
var ctx context.Context
var signer, _ = address.NewFromString("f3xdbznk6utswfpqclzzkxcamzshkwp2lwlsdq75oufyb3zfaaxncvy3qweuybo6elma4ypuolz7jptxk2m5ca")

func setup() {
	address.CurrentNetwork = address.Mainnet
	ctx = context.TODO()
	var err, endpoint, token = (error)(nil),
		"http://localhost:5678/rpc/v0",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.Hk_9_gbQoF-Z14yqIKa48h8Zmgbdy4WUMnGLRVjxLg4"
	header := http.Header{}
	header.Add("Authorization", "Bearer "+string(token))
	if walletAPI, closer, err = NewFullNodeRPC(context.TODO(), endpoint, header); err != nil {
		panic(err)
	}
	if signer, err = walletAPI.WalletNew(ctx, types.KTSecp256k1); err != nil {
		panic(err)
	}
}

func shutdown() {
	closer()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestAPI_WalletNew(t *testing.T) {
	addr, err := walletAPI.WalletNew(ctx, types.KTBLS)
	if err != nil {
		t.Errorf("wallet new failed:%s", err.Error())
	}
	t.Logf("wallet new bls:%s", addr.String())

	signer = addr

	addr, err = walletAPI.WalletNew(ctx, types.KTSecp256k1)
	if err != nil {
		t.Errorf("wallet new failed:%s", err.Error())
	}
	t.Logf("wallet new secp256k1:%s", addr.String())
}

func TestAPI_WalletHas(t *testing.T) {
	addr, _ := address.NewFromString("f3xdbznk6utswfpqclzzkxcamzshkwp2lwlsdq75oufyb3zfaaxncvy3qweuybo6elma4ypuolz7jptxk2m5ca")
	find, err := walletAPI.WalletHas(ctx, addr)
	if err != nil {
		t.Errorf("wallet has failed:%s", err.Error())
	}
	t.Logf("wallet has(%s) : %v", addr.String(), find)
}

type SignItem struct {
	Signer address.Address
	ToSign []byte
	Meta   api2.MsgMeta
}

func TestAPI_WalletSign(t *testing.T) {
	var tosign *SignItem
	var c, _ = cid.Decode("bafyreicmaj5hhoy5mgqvamfhgexxyergw7hdeshizghodwkjg6qmpoco7i")
	for mt, _ := range db_proc.SupportedMsgTypes {
		switch mt {
		case api2.MTChainMsg:
			msg := &types.Message{
				To:         builtin.StoragePowerActorAddr,
				From:       signer,
				Nonce:      0,
				Value:      big.Zero(),
				GasLimit:   81,
				GasFeeCap:  big.NewInt(234),
				GasPremium: big.NewInt(234),
				Method:     6,
				Params:     []byte("hai..."),
			}
			extal, _ := msg.Serialize()
			tosign = &SignItem{signer, msg.Cid().Bytes(), api2.MsgMeta{
				Type:  mt,
				Extra: extal}}
		case api2.MTDealProposal:
			// src code of signing dealproposal:
			// storagemarket/impl/client.go:372 -> markets/storageadapter/client.go:316
			proposal := shared_testutil.MakeTestUnsignedDealProposal()
			tosign = &SignItem{Signer: signer, Meta: api2.MsgMeta{mt, nil}, ToSign: unsafeCborUnmarshal(&proposal)}
		case api.MTClientDeal:
			// src code of signing clientdeal: storagemarket/impl/client.go:330
			inst := shared_testutil.MakeTestClientDealProposal()
			tosign = &SignItem{Signer: signer, Meta: api2.MsgMeta{mt, nil}, ToSign: unsafeCborUnmarshal(inst)}
		case api2.MTBlock:
			inst := &types.BlockHeader{
				Miner: signer,
				Ticket: &types.Ticket{
					VRFProof: []byte("vrf proof0000000vrf proof0000000"),
				},
				ElectionProof: &types.ElectionProof{
					VRFProof: []byte("vrf proof0000000vrf proof0000000"),
				},
				Parents:               []cid.Cid{c, c},
				ParentMessageReceipts: c,
				BLSAggregate:          &crypto.Signature{Type: crypto.SigTypeBLS, Data: []byte("boo! im a signature")},
				ParentWeight:          types.NewInt(123125126212),
				Messages:              c,
				Height:                85919298723,
				ParentStateRoot:       c,
				BlockSig:              &crypto.Signature{Type: crypto.SigTypeBLS, Data: []byte("boo! im a signature")},
				ParentBaseFee:         types.NewInt(3432432843291),
			}
			tosign = &SignItem{Signer: signer, Meta: api2.MsgMeta{mt, nil}, ToSign: unsafeCborUnmarshal(inst)}
		case api.MTDrawRandomParam:
			drp := &api.DrawRandomParams{Rbase: []byte("hello abc"), Pers: 10, Round: 200, Entropy: []byte("entry")}
			tosign = &SignItem{Signer: signer, Meta: api2.MsgMeta{mt, nil}, ToSign: unsafeCborUnmarshal(drp)}
		case api.MTSignedVoucher:
			inst := shared_testutil.MakeTestSignedVoucher()
			tosign = &SignItem{Signer: signer, Meta: api2.MsgMeta{mt, nil}, ToSign: unsafeCborUnmarshal(inst)}
		case api.MTStorageAsk:
			var ask *storagemarket.StorageAsk = shared_testutil.MakeTestStorageAsk()
			tosign = &SignItem{Signer: signer, Meta: api2.MsgMeta{mt, nil}, ToSign: unsafeCborUnmarshal(ask)}
		case api.MTAskResponse:
			var inst network.AskResponse = shared_testutil.MakeTestStorageAskResponse()
			tosign = &SignItem{Signer: signer, Meta: api2.MsgMeta{mt, nil}, ToSign: unsafeCborUnmarshal(&inst)}
		case api.MTNetWorkResponse:
			var resp network.Response = shared_testutil.MakeTestStorageNetworkResponse()
			tosign = &SignItem{Signer: signer, Meta: api2.MsgMeta{mt, nil}, ToSign: unsafeCborUnmarshal(&resp)}
		default:
			t.Errorf("unkown type:%s", mt)
		}
		signature, err := walletAPI.WalletSign(ctx, tosign.Signer, tosign.ToSign, tosign.Meta)
		if err != nil {
			t.Errorf("sign type(%s) failed:%s\n", tosign.Meta.Type, err.Error())
		} else {
			t.Logf("sign type success:%s, signature\n:%s\n",
				tosign.Meta.Type, toString(signature))
		}
	}
}

func TestAPI_WalletImportExportDel(t *testing.T) {
	var key, err, find, addrs = (*types.KeyInfo)(nil), error(nil), false, []address.Address(nil)
	var addr address.Address
	var retry = true
retry:
	if addrs, err = walletAPI.WalletList(ctx); err != nil {
		t.Errorf("wallet list failed:%s", err.Error())
	} else if len(addrs) == 0 {
		TestAPI_WalletNew(t)
		if retry {
			retry = false
			goto retry
		}
	} else {
		addr = addrs[0]
	}

	t.Logf("use address(%s) to test import, export , del", addr.String())

	if key, err = walletAPI.WalletExport(ctx, addr); err != nil {
		t.Errorf("wallet(%s) export failed:%s", addr.String(), err.Error())
		return
	}
	t.Logf("WalletExport(%s) success:\n%s\n", addr.String(), toString(key))

	if find, err = walletAPI.WalletHas(ctx, addr); err != nil {
		t.Errorf("WalletHas(%s) failed:%s", addr.String(), err.Error())
		return
	} else if !find {
		t.Errorf("WalletHas(%s) failed, address must exists!", addr.String())
		return
	}

	t.Logf("WalletHas(%s) success, address exists!", addr.String())

	if err = walletAPI.WalletDelete(ctx, addr); err != nil {
		t.Errorf("WalletDelete(%s) failed:%s",
			signer.String(), err.Error())
	}

	if find, err = walletAPI.WalletHas(ctx, addr); err != nil {
		t.Errorf("WalletHash(%s) failed:%s",
			addr.String(), err.Error())
		return
	} else if find {
		t.Errorf("WalletHas(%s) failed, address must not exists!", addr.String())
		return
	}

	t.Logf("HasWallet(%s) success, wallet had been deleted, doesn't exists!", addr.String())

	var expAddr address.Address
	if expAddr, err = walletAPI.WalletImport(ctx, key); err != nil {
		t.Errorf("WalletImport(%s) failed:%s", addr.String(), err.Error())
		return
	} else if expAddr != addr {
		t.Errorf("WalletImport(%s) failed, returned address(%s) not equals imported address",
			addr.String(), expAddr.String())
		return
	}

	t.Logf("WalletImport(%s) success", addr.String())

	if find, err = walletAPI.WalletHas(ctx, addr); err != nil {
		t.Errorf("WalletHash(%s) failed:%s",
			addr.String(), err.Error())
		return
	} else if !find {
		t.Errorf("WalletHas(%s) failed, wallet had just been imported!", addr.String())
		return
	}

	t.Logf("HasWallet(%s) success, wallet exists!", addr.String())

}

func toString(i interface{}) string {
	x, e := json.MarshalIndent(i, "", "  ")
	if e != nil {

		fmt.Printf("-------->\nwarns: interface to String failed:%s\n", e.Error())
	}
	return string(x)
}

func unsafeCborUnmarshal(in interface{}) []byte {
	bytes, err := cborutil.Dump(in)
	if err != nil {
		fmt.Printf("-------->\nwarns: fcborutils.Dump failed:%s\n", err.Error())
	}
	return bytes
}
