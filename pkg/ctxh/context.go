package ctxh

import "context"

type ContextValueHolder func(ctx context.Context) context.Context

func AttachContext(ctx context.Context, holders ...ContextValueHolder) context.Context {
	newCtx := ctx
	for _, holder := range holders {
		newCtx = holder(newCtx)
	}
	return newCtx
}
