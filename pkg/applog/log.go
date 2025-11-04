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

func Debug(ctx context.Context, msg string, args ...any) {
	loggerOrPanic(ctx).DebugContext(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	loggerOrPanic(ctx).InfoContext(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	loggerOrPanic(ctx).WarnContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	loggerOrPanic(ctx).ErrorContext(ctx, msg, args...)
}

func With(ctx context.Context, args ...any) context.Context {
	logger := loggerOrPanic(ctx)
	logger = logger.With(args...)
	return ctxval.Register(ctx, logger)
}
