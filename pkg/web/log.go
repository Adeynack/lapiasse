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

type logFormatter struct{}

var _ middleware.LogFormatter = (*logFormatter)(nil)

// NewLogEntry implements middleware.LogFormatter.
func (l *logFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &logEntry{
		ctx: r.Context(),
		reqArgs: []any{
			slog.String("endpoint", fmt.Sprintf(
				"%s %s://%s%s %s\" ",
				r.Method,
				lo.Ternary(r.TLS == nil, "http", "https"),
				r.Host,
				r.RequestURI,
				r.Proto,
			)),
			slog.String("remote_addr", r.RemoteAddr),
		},
	}
}

type logEntry struct {
	ctx     context.Context
	reqArgs []any
}

var _ middleware.LogEntry = (*logEntry)(nil)

func (l *logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	args := append(l.reqArgs,
		slog.Int("status", status),
		slog.Int("response_size_b", bytes),
		slog.Duration("duration", elapsed),
	)

	applog.Info(l.ctx, "Web request answered", args...)
}

func (l *logEntry) Panic(v any, stack []byte) {
	stackLines := strings.Split(string(stack), "\n")
	stackLines = lo.FilterMap(stackLines, func(l string, _ int) (string, bool) {
		l = strings.ReplaceAll(l, "\t", "  ")
		return l, l != ""
	})

	args := append(l.reqArgs,
		slog.Any("panic", v),
		slog.Any("stack", stackLines),
	)

	applog.Error(l.ctx, "Web request panicked", args...)
}
