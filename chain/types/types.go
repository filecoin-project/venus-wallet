package types

import (
	"github.com/filecoin-project/lotus/chain/types"
)

type SignedMsg struct {
	Cid       string
	Nonce     uint64
	SignedMsg types.Message
	Epoch     uint64
}
