package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"adeynack.net/lapiasse/pkg/applog"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/samber/lo"
)

type logFormatter struct {
	// Split logs into one upon request received and another upon response sent.
	// Otherwise (default), only one log entry is created per request, just before response is sent.
	SplitLogs bool
}

// NewLogEntry implements [middleware.LogFormatter.NewLogEntry].
func (l *logFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &logEntry{
		ctx:       r.Context(),
		SplitLogs: l.SplitLogs,
		sharedArgs: []any{
			slog.String("endpoint", fmt.Sprintf(
				"%s %s://%s%s %s\" ",
				r.Method,
				lo.Ternary(r.TLS == nil, "http", "https"),
				r.Host,
				r.RequestURI,
				r.Proto,
			)),
		},
		requestArgs: []any{
			slog.String("remote_addr", r.RemoteAddr),
		},
	}

	if entry.SplitLogs {
		applog.Info(r.Context(), "Web request received", append(entry.sharedArgs, entry.requestArgs...)...)
	}

	return entry
}

type logEntry struct {
	ctx       context.Context
	SplitLogs bool

	requestArgs []any
	sharedArgs  []any
}

// Write implements [middleware.LogEntry.Write].
func (l *logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra any) {
	// Log the shared args.
	args := l.sharedArgs

	// If in split mode, the request args are already logged.
	if !l.SplitLogs {
		args = append(args, l.requestArgs...)
	}

	// Add response info
	args = append(args,
		slog.Int("status", status),
		slog.Int("response_size_b", bytes),
		slog.Duration("duration", elapsed),
	)

	applog.Info(l.ctx, "Web request answered", args...)
}

// Panic implements [middleware.LogEntry.Panic].
func (l *logEntry) Panic(v any, stack []byte) {
	stackLines := strings.Split(string(stack), "\n")
	stackLines = lo.FilterMap(stackLines, func(l string, _ int) (string, bool) {
		l = strings.ReplaceAll(l, "\t", "  ")
		return l, l != ""
	})

	panicArgs := []any{
		slog.Any("panic", v),
		slog.Any("stack", stackLines),
	}

	applog.Error(l.ctx, "Web request panicked", append(l.sharedArgs, panicArgs...)...)
}
