//go:build test

package app

import (
	"context"
	"testing"

	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/controller"
	"adeynack.net/lapiasse/pkg/env"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"adeynack.net/lapiasse/pkg/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateTestAppCtx(t testing.TB) context.Context {
	ctx := t.Context()

	// Environments
	ctx = ctxval.RegisterNamed(ctx, "build", env.EnvTest)
	ctx = ctxval.RegisterNamed(ctx, "run", env.EnvTest)

	// The CleanupRecorder will ensure that all registered cleanup functions are called when the test ends.
	ctx = ctxval.Register(ctx, ctxval.CleanupRecorder(func(f ctxval.CleanupFunc) {
		t.Cleanup(func() {
			assert.NoError(t, f(ctx))
		})
	}))

	// Logger
	ctx = applog.RegisterTestLogger(ctx, t)

	// Database
	testDb, err := repository.InitializeGorm(ctx, &repository.Configuration{
		InMemory: true,
	})
	require.NoError(t, err)
	ctx = ctxval.Register(ctx, testDb)

	// API Implementation
	ctx = ctxval.Register(ctx, controller.New())

	return ctx
}
