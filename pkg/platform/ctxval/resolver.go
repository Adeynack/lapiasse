package ctxval

import (
	"context"
	"time"
)

// Resolver offers a facade to hash-based dependency resolution, offering
// an alternative to the recursive calls of the basic [context.Context] implementation.
//
// It implements the [context.Context] interface, so it can be used as a drop-in
// replacement for contexts, while providing faster dependency resolution capabilities.
type Resolver struct {
	ctx               context.Context
	dependenciesByKey map[contextValueKey]any
}

// NewResolver creates a new Resolver instance, wrapping the given context.
func NewResolver(ctx context.Context) *Resolver {
	if ctx == nil {
		ctx = context.Background()
	}

	return &Resolver{
		ctx:               ctx,
		dependenciesByKey: make(map[contextValueKey]any),
	}
}

// Deadline implements the [context.Context] interface.
func (r *Resolver) Deadline() (deadline time.Time, ok bool) {
	return r.ctx.Deadline()
}

// Done implements the [context.Context] interface.
func (r *Resolver) Done() <-chan struct{} {
	return r.ctx.Done()
}

// Err implements the [context.Context] interface.
func (r *Resolver) Err() error {
	return r.ctx.Err()
}

// Value implements the [context.Context] interface.
func (r *Resolver) Value(key any) any {
	switch key := key.(type) {
	case contextValueKey:
		if v := r.dependenciesByKey[key]; v != nil {
			return v
		}
	}

	return r.ctx.Value(key)
}

// RegisterInResolver registers a dependency of type T in the resolver,
// to be retrieved later by type.
func RegisterInResolver[T any](resolver *Resolver, value T) {
	key := keyFor[T]("")
	resolver.dependenciesByKey[key] = value
}

// RegisterNamedInResolver registers a named dependency of type T in the resolver,
// to be retrieved later by type and name.
func RegisterNamedInResolver[T any](resolver *Resolver, name string, value T) {
	key := keyFor[T](name)
	resolver.dependenciesByKey[key] = value
}
