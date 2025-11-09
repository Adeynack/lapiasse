//go:build test

package applog

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/golang-cz/devslog"
)

func RegisterTestLogger(ctx context.Context, t testing.TB) context.Context {
	// // for now, just register the default logger; otherwise the output is quite poluted.
	// return WithLogger(ctx, slog.Default())

	buffer := new(bytes.Buffer)
	handler := devslog.NewHandler(buffer, &devslog.Options{
		HandlerOptions: &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
		NewLineAfterLog: true,
	})
	logger := slog.New(handler)

	// Ensure that at the end of the test, the buffer is printed to t.Logf if the test failed.
	t.Cleanup(func() {
		if t.Failed() {
			t.Logf("Logger output:\n%s", buffer.String())
		}
	})

	return WithLogger(ctx, logger)
}
