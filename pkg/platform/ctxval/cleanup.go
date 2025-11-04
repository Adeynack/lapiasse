package ctxval

import (
	"context"
)

// CleanupFunc defines a function that performs cleanup operations using the provided context.
type CleanupFunc func(context.Context)

// CleanupRecorder defines a function that records cleanup functions to be executed later.
// It offers a type to register in a context to allow other values to register their own cleanup functions.
type CleanupRecorder func(CleanupFunc)
