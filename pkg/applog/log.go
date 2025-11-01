package applog

import "context"

func Info(ctx context.Context, msg string, args ...any) {
	loggerOrPanic(ctx).InfoContext(ctx, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	loggerOrPanic(ctx).DebugContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	loggerOrPanic(ctx).ErrorContext(ctx, msg, args...)
}
