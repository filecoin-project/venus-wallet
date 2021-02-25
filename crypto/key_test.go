package crypto

import (
	"encoding/hex"
	"fmt"
	"github.com/ipfs-force-community/venus-wallet/core"
	"gotest.tools/assert"
	"testing"
)

func TestSecpPrivateKey(t *testing.T) {
	k := "2a2a2a2a2a2a2a2a5fbf0ed0f8364c01ff27540ecd6669ff4cc548cbe60ef5ab"
	key, err := hex.DecodeString(k)
	if err != nil {
		t.Fatal(err)
	}
	prv, err := NewKeyFromData2(core.KTSecp256k1, key)
	if err != nil {
		t.Fatalf("load private key from bytes err:%s", err)
	}
	pub := prv.Public()
	assert.Equal(t, fmt.Sprintf("%x", pub), "047a7735f3935272bda6748144b16fc7c036137971af61f24ec2a197ac1e8eb237dd4f3d24f9f88b1e6762ad829c1b68386506cd7cb44ade7e484b19740690d4bf")

	addr, err := prv.Address()
	if err != nil {
		t.Fatalf("private key parse address error:%s", err)
	}
	assert.Equal(t, addr.String(), "t17rmoshisxovfjkrox2gpryaupmsalzy2tlaluiq")

	signData := []byte("hello filecoin")
	signature, err := prv.Sign(signData)
	if err != nil {
		t.Fatalf("private key sign data err:%s", err)
	}
	assert.Equal(t, fmt.Sprintf("%x", signature.Data), "4c49bacbd5a1e9734595d77a7ae8909ae35859e182344c0f1c081d8fdc6749302d50ae81d0d80c9c4be9cf3071bf7b0c465a19aea1023b21239b8bbdffd0ca8d01")
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
	prv, err := NewKeyFromData2(core.KTBLS, key)
	if err != nil {
		t.Fatalf("load private key from bytes err:%s", err)
	}

	pub := prv.Public()
	assert.Equal(t, fmt.Sprintf("%x", pub), "99ddf98471b6ef4c7b01cf9fd9945f6ca10069d93ca2d61ed86c43b51080262c4d948b40c2f3571ab0aea9b2e22656fc")

	addr, err := prv.Address()
	if err != nil {
		t.Fatalf("private key parse address error:%s", err)
	}
	assert.Equal(t, addr.String(), "t3tho7tbdrw3xuy6ybz6p5tfc7nsqqa2ozhsrnmhwynrb3keeaeywe3felidbpgvy2wcxktmxcezlpz7pbu5pq")

	signData := []byte("hello filecoin")
	signature, err := prv.Sign(signData)
	if err != nil {
		t.Fatalf("private key sign data err:%s", err)
	}
	assert.Equal(t, fmt.Sprintf("%x", signature.Data), "828e485cf906ae5deda33ec3e16b2a5ac761ab5a6109f481134a4dd40d7b5c3b03d187f614c7dd85d142ff50b902bdf50923f58b5bb741b4a522bf17f8efe5cf3f8f51715cd31a765eb5e5d76889102752d2cc2efa4287829cacb4de2be07cf7")
}
