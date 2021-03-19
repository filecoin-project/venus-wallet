package core

import "context"

type StrategyLevel = uint8

const (
	SLDisable StrategyLevel = iota
	SLMetaType
	SLMethod
)

var WalletStrategyLevel StrategyLevel

type ContextKey string

const CtxKeyStrategy ContextKey = "strategy"

func ContextStrategyToken(ctx context.Context) string {
	localPerm, ok := ctx.Value(CtxKeyStrategy).(string)
	if !ok {
		return StringEmpty
	}
	return localPerm
}
