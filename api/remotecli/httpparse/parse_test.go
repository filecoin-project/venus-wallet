package httpparse

import (
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestParseApiInfo(t *testing.T) {
	apiInfo, err := ParseApiInfo("/ip4/0.0.0.0/tcp/5678/http:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.L1IREykf3F6qqsyhSPaAHurHIrj2yz0ne04DyQ-YF-U")
	if err != nil {
		t.Fatal(err)
	}
	str2 := "/ip4/0.0.0.0/tcp/5678/http:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.L1IREykf3F6qqsyhSPaAHurHIrj2yz0ne04DyQ-YF-U:62d3c94c-86d1-11eb-b252-acde48001122"
	apiInfo2, err := ParseApiInfo(str2)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, apiInfo.Addr, apiInfo2.Addr)
	assert.Equal(t, strings.HasPrefix(string(apiInfo2.Token), string(apiInfo.Token)), true)
	assert.Equal(t, strings.HasSuffix(string(apiInfo2.Token), "62d3c94c-86d1-11eb-b252-acde48001122"), true)
}
