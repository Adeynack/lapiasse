package slogex

import (
	"log/slog"
	"testing"
)

func TestMultiplexHandler(t *testing.T) {
	t.Run("Implements slog.Handler", func(t *testing.T) {
		var _ slog.Handler = (*MultiplexHandlerConfig)(nil)
	})
}
