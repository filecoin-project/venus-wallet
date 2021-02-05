package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	rpcmetrics "github.com/filecoin-project/go-jsonrpc/metrics"
)

var (
	Version, _ = tag.NewKey("version")
	Commit, _  = tag.NewKey("commit")
)

var (
	LotusInfo       = stats.Int64("info", "Arbitrary counter to tag lotus info to", stats.UnitDimensionless)
	ChainNodeHeight = stats.Int64("chain/node_height", "Current Height of the node", stats.UnitDimensionless)
)

var (
	InfoView = &view.View{
		Name:        "info",
		Description: "Lotus node information",
		Measure:     LotusInfo,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{Version, Commit},
	}
	ChainNodeHeightView = &view.View{
		Measure:     ChainNodeHeight,
		Aggregation: view.LastValue(),
	}
)

// DefaultViews is an array of OpenCensus views for metric gathering purposes
var DefaultViews = append([]*view.View{
	InfoView,
	ChainNodeHeightView,
}, rpcmetrics.DefaultViews...)
