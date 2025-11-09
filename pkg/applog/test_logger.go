//go:build test

package applog

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/golang-cz/devslog"
)

func RegisterTestLogger(ctx context.Context, t testing.TB) context.Context {
	t.Helper()

	var testOutput io.Writer

	switch strings.ToLower(os.Getenv("TEST_LOG")) {
	case "off", "0", "":
		// disable logging completely during testing
		return WithLogger(ctx, slog.New(slog.DiscardHandler))
	case "all":
		// log all to the test's output
		testOutput = t.Output()
	case "fail":
		// log only if the test fails, to the test's output
		buffer := new(bytes.Buffer)
		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("Logger output:\n%s", buffer.String())
			}
		})
		testOutput = buffer
	}

	handler := devslog.NewHandler(testOutput, &devslog.Options{
		HandlerOptions: &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
		NewLineAfterLog: true,
	})
	logger := slog.New(handler)

	return WithLogger(ctx, logger)
}
