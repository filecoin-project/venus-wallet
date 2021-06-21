package node

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/asaskevich/EventBus"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/venus-wallet/api/remotecli/httpparse"
	"github.com/filecoin-project/venus-wallet/config"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/venus-wallet/core"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("node")

type TipSetKey struct {
	// The internal representation is a concatenation of the bytes of the CIDs, which are
	// self-describing, wrapped as a string.
	// These gymnastics make the a TipSetKey usable as a map key.
	// The empty key has value "".
	value string
}

// The length of a block header CID in bytes.
var blockHeaderCIDLen int

func init() {
	// hash a large string of zeros so we don't estimate based on inlined CIDs.
	var buf [256]byte
	c, err := abi.CidBuilder.Sum(buf[:])
	if err != nil {
		panic(err)
	}
	blockHeaderCIDLen = len(c.Bytes())
}

func (k TipSetKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.Cids())
}
func (k TipSetKey) Cids() []core.Cid {
	cids, err := decodeKey([]byte(k.value))
	if err != nil {
		panic("invalid tipset key: " + err.Error())
	}
	return cids
}
func decodeKey(encoded []byte) ([]cid.Cid, error) {
	// To avoid reallocation of the underlying array, estimate the number of CIDs to be extracted
	// by dividing the encoded length by the expected CID length.
	estimatedCount := len(encoded) / blockHeaderCIDLen
	cids := make([]cid.Cid, 0, estimatedCount)
	nextIdx := 0
	for nextIdx < len(encoded) {
		nr, c, err := cid.CidFromBytes(encoded[nextIdx:])
		if err != nil {
			return nil, err
		}
		cids = append(cids, c)
		nextIdx += nr
	}
	return cids, nil
}

func (k *TipSetKey) UnmarshalJSON(b []byte) error {
	var cids []core.Cid
	if err := json.Unmarshal(b, &cids); err != nil {
		return err
	}
	k.value = string(encodeKey(cids))
	return nil
}

func encodeKey(cids []core.Cid) []byte {
	buffer := new(bytes.Buffer)
	for _, c := range cids {
		// bytes.Buffer.Write() err is documented to be always nil.
		_, _ = buffer.Write(c.Bytes())
	}
	return buffer.Bytes()
}

type Actor struct {
	// Identifies the type of actor (string coded as a CID), see `chain/actors/actors.go`.
	Code    core.Cid
	Head    core.Cid
	Nonce   uint64
	Balance big.Int
}

// NodeClient connect Lotus or Venus node and call json RPC API
type NodeClient struct {
	// NOTE:
	StateGetActor func(ctx context.Context, actor core.Address, tsk TipSetKey) (*Actor, error)
	Cancel        func()
}

var EmptyNodeClient = &NodeClient{}

func NewNodeClient(cnf *config.StrategyConfig) (*NodeClient, error) {
	if cnf.Level < core.SLMethod {
		return EmptyNodeClient, nil
	}
	if cnf.NodeURL == core.StringEmpty {
		return nil, errors.New("node url can not be empty when level is SLMethod")
	}
	ai, err := httpparse.ParseApiInfo(cnf.NodeURL)
	if err != nil {
		return nil, err
	}
	addr, err := ai.DialArgs()
	if err != nil {
		return nil, err
	}
	cli := &NodeClient{}
	closer, err := jsonrpc.NewClient(context.Background(), addr, "Filecoin", cli, ai.AuthHeader())
	if err != nil {
		return nil, err
	}
	cli.Cancel = closer
	log.Info("node client initialize successfully")
	return cli, nil
}

func NewEventBus() EventBus.Bus {
	return EventBus.New()
}
