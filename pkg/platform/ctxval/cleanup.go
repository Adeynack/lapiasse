package ctxval

import (
	"context"
)

// CleanupFunc defines a function that performs cleanup operations using the provided context.
type CleanupFunc func(context.Context) error

// CleanupRecorder defines a function that records cleanup functions to be executed later.
// It offers a type to register in a context to allow other values to register their own cleanup functions.
type CleanupRecorder func(CleanupFunc)

// MustCleanup registers a cleanup function in the context's [CleanupRecorder].
// It panics if the [CleanupRecorder] cannot be resolved from the context.
func MustCleanup(ctx context.Context, f CleanupFunc) {
	recorder := MustResolve[CleanupRecorder](ctx)
	recorder(f)
}
