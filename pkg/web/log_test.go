package web

import (
	"testing"

	"github.com/go-chi/chi/v5/middleware"
)

func TestLogFormatter(t *testing.T) {
	t.Run("implements middleware.LogFormatter", func(t *testing.T) {
		var _ middleware.LogFormatter = (*logFormatter)(nil)
	})
}

func TestLogEntry(t *testing.T) {
	t.Run("implements middleware.LogEntry", func(t *testing.T) {
		var _ middleware.LogEntry = (*logEntry)(nil)
	})
}
