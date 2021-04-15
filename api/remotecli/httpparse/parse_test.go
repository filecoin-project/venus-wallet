package httpparse

import (
	"testing"

	"gotest.tools/assert"
)

func TestParseApiInfo(t *testing.T) {
	apiInfo, err := ParseApiInfo("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.L1IREykf3F6qqsyhSPaAHurHIrj2yz0ne04DyQ-YF-U:/ip4/0.0.0.0/tcp/5678/http")
	if err != nil {
		t.Fatal(err)
	}
	str2 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.L1IREykf3F6qqsyhSPaAHurHIrj2yz0ne04DyQ-YF-U:/ip4/0.0.0.0/tcp/5678/http:62d3c94c-86d1-11eb-b252-acde48001122"
	apiInfo2, err := ParseApiInfo(str2)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, apiInfo.Addr.String(), apiInfo2.Addr.String())
	assert.DeepEqual(t, apiInfo.Token, apiInfo2.Token)
	assert.DeepEqual(t, apiInfo2.StrategyToken, []byte("62d3c94c-86d1-11eb-b252-acde48001122"))

	str3 := "62d3c94c-86d1-11eb-b252-acde48001122:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.L1IREykf3F6qqsyhSPaAHurHIrj2yz0ne04DyQ-YF-U:/ip4/0.0.0.0/tcp/5678/http"
	apiInfo3, err := ParseApiInfo(str3)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, apiInfo.Addr.String(), apiInfo3.Addr.String())
	assert.DeepEqual(t, apiInfo.Token, apiInfo3.Token)
	assert.DeepEqual(t, apiInfo3.StrategyToken, []byte("62d3c94c-86d1-11eb-b252-acde48001122"))

	ipv4 := "/ip4/0.0.0.0/tcp/5678/http"
	ipv4Info, err := ParseApiInfo(ipv4)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, ipv4Info.Addr.String(), ipv4)
	ipv4 = "/ip4/192.168.1.151/tcp/5678/http"
	ipv4Info, err = ParseApiInfo(ipv4)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, ipv4Info.Addr.String(), ipv4)
}
