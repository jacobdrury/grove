package config

import "context"

type contextKey string

const (
	noHooksContextKey = contextKey("noHooks")
	pipeContextKey    = contextKey("pipe")
)

func ContextWithNoHooks(ctx context.Context) context.Context {
	return context.WithValue(ctx, noHooksContextKey, true)
}

func NoHooks(ctx context.Context) bool {
	if value, ok := ctx.Value(noHooksContextKey).(bool); ok {
		return value
	}

	return false
}

func ContextWithPipe(ctx context.Context) context.Context {
	return context.WithValue(ctx, pipeContextKey, true)
}

func Pipe(ctx context.Context) bool {
	if value, ok := ctx.Value(pipeContextKey).(bool); ok {
		return value
	}

	return false
}
