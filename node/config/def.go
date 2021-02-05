package config

import (
	"encoding"
	"time"
)

// Common is common config between full node and miner
type Common struct {
	API       API
	DbCfg     DbCfg
}

// FullNode is a full node config
type FullNode struct {
	Common
	Metrics Metrics
}

// // Common

// API contains configs for API endpoint
type API struct {
	ListenAddress       string
	RemoteListenAddress string
	Timeout             Duration
}

// db config
type DbCfg struct {
	Conn      string `json:"conn" binding:"required"`
	Type      string `json:"type" binding:"required"`
	DebugMode bool   `json:"debugMode" binding:"required"`
}

// // Full Node
type Metrics struct {
	Nickname   string
	HeadNotifs bool
}

func defCommon() Common {
	return Common{
		API: API{
			ListenAddress: "/ip4/0.0.0.0/tcp/5678/http",
			Timeout:       Duration(30 * time.Second),
		},
		DbCfg: DbCfg{
			Conn:      "",
			Type:      "sqlite",
			DebugMode: true,
		},
	}

}

// DefaultFullNode returns the default config
func DefaultFullNode() *FullNode {
	return &FullNode{
		Common: defCommon(),
	}
}

var _ encoding.TextMarshaler = (*Duration)(nil)
var _ encoding.TextUnmarshaler = (*Duration)(nil)

// Duration is a wrapper type for time.Duration
// for decoding and encoding from/to TOML
type Duration time.Duration

// UnmarshalText implements interface for TOML decoding
func (dur *Duration) UnmarshalText(text []byte) error {
	d, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*dur = Duration(d)
	return err
}

func (dur Duration) MarshalText() ([]byte, error) {
	d := time.Duration(dur)
	return []byte(d.String()), nil
}
