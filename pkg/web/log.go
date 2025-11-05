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
	endpointArg := slog.String("endpoint", fmt.Sprintf(
		"%s %s://%s%s %s\" ",
		r.Method,
		lo.Ternary(r.TLS == nil, "http", "https"),
		r.Host,
		r.RequestURI,
		r.Proto,
	))

	applog.Info(r.Context(), "Web request received",
		endpointArg,
		slog.String("remote_addr", r.RemoteAddr),
	)

	return &logEntry{
		ctx:         r.Context(),
		endpointArg: endpointArg,
	}
}

type logEntry struct {
	ctx         context.Context
	endpointArg slog.Attr
}

var _ middleware.LogEntry = (*logEntry)(nil)

func (l *logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	applog.Info(l.ctx, "Web request answered",
		l.endpointArg,
		slog.Int("status", status),
		slog.Int("response_size_b", bytes),
		slog.Duration("duration", elapsed),
	)
}

func (l *logEntry) Panic(v any, stack []byte) {
	stackLines := strings.Split(string(stack), "\n")
	stackLines = lo.FilterMap(stackLines, func(l string, _ int) (string, bool) {
		l = strings.ReplaceAll(l, "\t", "  ")
		return l, l != ""
	})

	applog.Error(l.ctx, "Web request panicked",
		l.endpointArg,
		slog.Any("panic", v),
		slog.Any("stack", stackLines),
	)
}
