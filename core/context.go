package core

import "context"

var ContextEnable bool

type ContextKey string

const CtxKeyStrategy ContextKey = "strategy"

func ContextStrategyToken(ctx context.Context) string {
	localPerm, ok := ctx.Value(CtxKeyStrategy).(string)
	if !ok {
		return StringEmpty
	}
	return localPerm
}
