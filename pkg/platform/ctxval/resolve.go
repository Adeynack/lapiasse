package ctxval

import (
	"context"
	"fmt"
)

// Resolve finds in the provided ctx an unnamed resource of type T.
func Resolve[T any](ctx context.Context) (T, error) {
	return ResolveNamed[T](ctx, "")
}

// ResolveNamed finds in the provided ctx a resource of type T with the given name.
func ResolveNamed[T any](ctx context.Context, name string) (result T, err error) {
	key := keyFor[T](name)

	rawValue := ctx.Value(key)
	if rawValue == nil {
		return result, fmt.Errorf("%w %q", ErrUnregisteredDependency, key)
	}

	result, ok := rawValue.(T)
	if !ok {
		return result, fmt.Errorf("%w %q", ErrUnexpectedType, key)
	}

	return result, nil
}

// Resolve finds in the provided ctx an unnamed resource of type T or panics.
func MustResolve[T any](ctx context.Context) T {
	return MustResolveNamed[T](ctx, "")
}

// ResolveNamed finds in the provided ctx a resource of type T with the given name or panics.
func MustResolveNamed[T any](ctx context.Context, name string) T {
	result, err := ResolveNamed[T](ctx, name)
	if err != nil {
		panic(err)
	}

	return result
}

// Resolve2 is a shortcut to multiple calls to Resolve and returns the first encountered error.
func Resolve2[T1, T2 any](ctx context.Context) (val1 T1, val2 T2, err error) {
	val1, err = Resolve[T1](ctx)
	if err != nil {
		return val1, val2, err
	}

	val2, err = Resolve[T2](ctx)

	return val1, val2, err
}

// Resolve3 is a shortcut to multiple calls to Resolve and returns the first encountered error.
func Resolve3[T1, T2, T3 any](ctx context.Context) (val1 T1, val2 T2, val3 T3, err error) {
	val1, val2, err = Resolve2[T1, T2](ctx)
	if err != nil {
		return val1, val2, val3, err
	}

	val3, err = Resolve[T3](ctx)

	return val1, val2, val3, err
}

// Resolve4 is a shortcut to multiple calls to Resolve and returns the first encountered error.
func Resolve4[T1, T2, T3, T4 any](ctx context.Context) (val1 T1, val2 T2, val3 T3, val4 T4, err error) {
	val1, val2, val3, err = Resolve3[T1, T2, T3](ctx)
	if err != nil {
		return val1, val2, val3, val4, err
	}

	val4, err = Resolve[T4](ctx)

	return val1, val2, val3, val4, err
}
