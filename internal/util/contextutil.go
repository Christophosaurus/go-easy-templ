package util

import (
	"context"
)

func SetContext(ctx context.Context, key string, value any) context.Context {
	ctx = context.WithValue(ctx, key, value)
	return ctx
}

func GetContextString(ctx context.Context, key string) (string, bool) {
	value, ok := ctx.Value(key).(string)
	return value, ok
}
