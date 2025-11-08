//go:build test

package app

import (
	"context"
	"log/slog"
	"testing"

	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"adeynack.net/lapiasse/pkg/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateTestAppCtx(t testing.TB) context.Context {
	ctx := t.Context()

	// The CleanupRecorder will ensure that all registered cleanup functions are called when the test ends.
	ctx = ctxval.Register(ctx, ctxval.CleanupRecorder(func(f ctxval.CleanupFunc) {
		t.Cleanup(func() {
			assert.NoError(t, f(ctx))
		})
	}))

	// todo: Make sure that the default logger during tests is logging to the test logger (t.Logf)
	ctx = ctxval.Register(ctx, slog.Default())

	testDb, err := repository.InitializeGorm(ctx, &repository.Configuration{
		InMemory: true,
	})
	require.NoError(t, err)
	ctx = ctxval.Register(ctx, testDb)

	return ctx
}
