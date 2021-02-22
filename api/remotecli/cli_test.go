package remotecli

import "testing"

func TestParseApiInfo(t *testing.T) {
	apiInfo, err := ParseApiInfo("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.L1IREykf3F6qqsyhSPaAHurHIrj2yz0ne04DyQ-YF-U:/ip4/0.0.0.0/tcp/5678/http")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(apiInfo)
}
