package ctxval

import (
	"context"
)

// Container offers a facade to hash-based dependency resolution, offering
// an alternative to the recursive calls of the basic [context.Context] implementation.
//
// It implements the [context.Context] interface, so it can be used as a drop-in
// replacement for contexts, while providing faster dependency resolution capabilities.
type Container struct {
	context.Context
	dependenciesByKey map[contextValueKey]any
}

// NewContainer creates a new Container instance, wrapping the given context.
func NewContainer(parent context.Context) *Container {
	if parent == nil {
		panic("cannot create container from nil parent context")
	}

	return &Container{
		Context:           parent,
		dependenciesByKey: make(map[contextValueKey]any),
	}
}

// Value implements the [context.Context] interface.
//
// It first checks the container's internal dependency map,
// before falling back to the wrapped context's Value method.
func (r *Container) Value(key any) any {
	switch key := key.(type) {
	case contextValueKey:
		if v := r.dependenciesByKey[key]; v != nil {
			return v
		}
	}

	return r.Context.Value(key)
}

// RegisterInContainer registers a dependency of type T in the container,
// to be retrieved later by type.
func RegisterInContainer[T any](container *Container, value T) {
	key := keyFor[T]("")
	container.dependenciesByKey[key] = value
}

// RegisterNamedInContainer registers a named dependency of type T in the container,
// to be retrieved later by type and name.
func RegisterNamedInContainer[T any](container *Container, name string, value T) {
	key := keyFor[T](name)
	container.dependenciesByKey[key] = value
}
