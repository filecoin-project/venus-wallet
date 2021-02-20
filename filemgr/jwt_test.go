package filemgr

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestRandSecret(t *testing.T) {
	sec, err := randSecret()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(sec.key))

}
