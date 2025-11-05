package slogex

import (
	"context"
	"log/slog"
)

type MultiplexHandlerConfig struct {
	handlers []slog.Handler
}

// NewMultiplexHandler creates a new [slog.Handler] that forwards log records
// to multiple underlying handlers.
func NewMultiplexHandler(handlers ...slog.Handler) slog.Handler {
	return &MultiplexHandlerConfig{
		handlers: handlers,
	}
}

// Enabled implements [slog.Handler.Enabled].
func (m *MultiplexHandlerConfig) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}

	return false
}

// Handle implements [slog.Handler.Handle].
func (m *MultiplexHandlerConfig) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if err := h.Handle(ctx, r); err != nil {
			return err
		}
	}

	return nil
}

// WithAttrs implements [slog.Handler.WithAttrs].
func (m *MultiplexHandlerConfig) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}

	return &MultiplexHandlerConfig{
		handlers: newHandlers,
	}
}

// WithGroup implements [slog.Handler.WithGroup].
func (m *MultiplexHandlerConfig) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}

	return &MultiplexHandlerConfig{
		handlers: newHandlers,
	}
}
