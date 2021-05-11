package core

import "context"

type StrategyLevel = uint8

const (
	// default level, no strategy
	SLDisable StrategyLevel = iota
	// level 1, check struct type
	SLMetaType
	// level 2, check struct type and method in msg
	SLMethod
)

// wallet strategy level, it will be initialized from config at startup
var WalletStrategyLevel StrategyLevel

type ContextKey string

// ctx key of strategy token
const CtxKeyStrategy ContextKey = "strategy"

func ContextStrategyToken(ctx context.Context) string {
	localPerm, ok := ctx.Value(CtxKeyStrategy).(string)
	if !ok {
		return StringEmpty
	}
	return localPerm
}
