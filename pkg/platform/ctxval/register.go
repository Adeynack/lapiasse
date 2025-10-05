package ctxval

import "context"

func Register[T any](ctx context.Context, value T) context.Context {
	return RegisterNamed(ctx, "", value)
}

func RegisterNamed[T any](ctx context.Context, name string, value T) context.Context {
	key := keyFor[T](name)
	return context.WithValue(ctx, key, value)
}
