//go:build test

package applog

import (
	"context"
	"log/slog"
	"testing"

	"github.com/golang-cz/devslog"
)

func RegisterTestLogger(ctx context.Context, t testing.TB) context.Context {
	handler := devslog.NewHandler(t.Output(), &devslog.Options{
		HandlerOptions: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
		NewLineAfterLog: true,
	})

	logger := slog.New(handler)

	return WithLogger(ctx, logger)
}
