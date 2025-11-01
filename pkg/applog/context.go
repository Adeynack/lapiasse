package applog

import (
	"context"
	"log/slog"

	"adeynack.net/lapiasse/pkg/platform/ctxval"
)

func loggerOrPanic(ctx context.Context) *slog.Logger {
	logger, err := ctxval.Resolve[*slog.Logger](ctx)
	if err != nil {
		panic(err)
	}

	return logger
}
