package crypto

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/filecoin-project/venus/venus-shared/types"
	"gotest.tools/assert"
)

func TestSecpPrivateKey(t *testing.T) {
	k := "2a2a2a2a2a2a2a2a5fbf0ed0f8364c01ff27540ecd6669ff4cc548cbe60ef5ab"
	key, err := hex.DecodeString(k)
	if err != nil {
		t.Fatal(err)
	}
	prv, err := NewKeyFromData2(types.KTSecp256k1, key)
	if err != nil {
		t.Fatalf("load private key from bytes err:%s", err)
	}
	pub := prv.Public()
	assert.Equal(t, fmt.Sprintf("%x", pub), "047a7735f3935272bda6748144b16fc7c036137971af61f24ec2a197ac1e8eb237dd4f3d24f9f88b1e6762ad829c1b68386506cd7cb44ade7e484b19740690d4bf")

	addr, err := prv.Address()
	if err != nil {
		t.Fatalf("private key parse address error:%s", err)
	}
	assert.Equal(t, addr.String(), "f17rmoshisxovfjkrox2gpryaupmsalzy2tlaluiq")

	signData := []byte("hello filecoin")
	signature, err := prv.Sign(signData)
	if err != nil {
		t.Fatalf("private key sign data err:%s", err)
	}
	assert.Equal(t, fmt.Sprintf("%x", signature.Data), "34a3d789cee65244de625a147d90af6c57cc1875fae8077689c5364298208e931b38122203a47e89990ab79a059759783bf57b224f27b63d632c5c7fea2da04c00")
}

func TestBls2KGen(t *testing.T) {
	key := "7b22707269766174654b6579223a22394d35307259715a31335673747a41557676576b5144594a323831724f32343557394e487a586d5448574d3d222c2274797065223a22626c73227d"
	prv := []byte(key)
	t.Logf("%x", prv)
}

func TestBLSPrivateKey(t *testing.T) {
	k := "33ff38fa8c1c53c3ef1b2f811cee9ce4f7ec2ce90b16ceeb9675dc6aa3d04817"
	key, err := hex.DecodeString(k)
	if err != nil {
		t.Fatal(err)
	}
	prv, err := NewKeyFromData2(types.KTBLS, key)
	if err != nil {
		t.Fatalf("load private key from bytes err:%s", err)
	}

	pub := prv.Public()
	assert.Equal(t, fmt.Sprintf("%x", pub), "99ddf98471b6ef4c7b01cf9fd9945f6ca10069d93ca2d61ed86c43b51080262c4d948b40c2f3571ab0aea9b2e22656fc")

	addr, err := prv.Address()
	if err != nil {
		t.Fatalf("private key parse address error:%s", err)
	}
	assert.Equal(t, addr.String(), "f3tho7tbdrw3xuy6ybz6p5tfc7nsqqa2ozhsrnmhwynrb3keeaeywe3felidbpgvy2wcxktmxcezlpz7pbu5pq")

	signData := []byte("hello filecoin")
	signature, err := prv.Sign(signData)
	if err != nil {
		t.Fatalf("private key sign data err:%s", err)
	}
	assert.Equal(t, fmt.Sprintf("%x", signature.Data), "828e485cf906ae5deda33ec3e16b2a5ac761ab5a6109f481134a4dd40d7b5c3b03d187f614c7dd85d142ff50b902bdf50923f58b5bb741b4a522bf17f8efe5cf3f8f51715cd31a765eb5e5d76889102752d2cc2efa4287829cacb4de2be07cf7")
}

func TestDelegatedPrivateKey(t *testing.T) {
	k := "33ff38fa8c1c53c3ef1b2f811cee9ce4f7ec2ce90b16ceeb9675dc6aa3d04821"
	key, err := hex.DecodeString(k)
	if err != nil {
		t.Fatal(err)
	}
	prv, err := NewKeyFromData2(types.KTDelegated, key)
	if err != nil {
		t.Fatalf("load private key from bytes err:%s", err)
	}

	pub := prv.Public()
	assert.Equal(t, fmt.Sprintf("%x", pub), "04f047588679644c440d8f7878fe29730e0cdc8c21b32ab8be5c8c01e1bf5d25e09609fff9c037fcd5cd6ac090e44a438c3d13a31a39a96a4b412831e4a3c83b0e")

	addr, err := prv.Address()
	if err != nil {
		t.Fatalf("private key parse address error:%s", err)
	}
	assert.Equal(t, addr.String(), "f410fnz6o7a4owcgw33z2zvsxtuf64aoxlmcn4dh5hzy")

	signData := []byte("hello filecoin")
	signature, err := prv.Sign(signData)
	if err != nil {
		t.Fatalf("private key sign data err:%s", err)
	}
	assert.Equal(t, fmt.Sprintf("%x", signature.Data), "48afd9f624ea4159c86e4c04fa517ec8e6e1e403a6fafda38df9499c54b505aa77504b3100d07e894490e039d39141d887f27b843ca43bf893ceb0c78f0034e700")

	assert.NilError(t, delegatedVerify(signature.Data, addr, signData))
}
