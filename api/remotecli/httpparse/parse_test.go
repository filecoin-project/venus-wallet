package httpparse

import (
	"testing"
)

func TestParseApiInfo(t *testing.T) {
	apiInfo, err := ParseApiInfo("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.L1IREykf3F6qqsyhSPaAHurHIrj2yz0ne04DyQ-YF-U:/ip4/0.0.0.0/tcp/5678/http")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(apiInfo)
	str2 :="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIl19.t3Bs0O7Tw89JRW4w-P87WXT4sqBIFtAYgMHdDYvuKp4:/ip4/0.0.0.0/tcp/5678/http:62d3c94c-86d1-11eb-b252-acde48001122"
	ai2 ,err:= ParseApiInfo(str2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(ai2.Token))
	t.Log(string(ai2.StrategyToken))
	t.Log(ai2.Addr.String())
}


