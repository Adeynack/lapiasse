package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"

	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/env"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"adeynack.net/lapiasse/pkg/repository"
)

func configureLogger(ctx context.Context, config *repository.Configuration) (*slog.Logger, error) {
	cleanup, err := ctxval.Resolve[ctxval.CleanupRecorder](ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving cleanup recorder: %w", err)
	}

	logFilePath := path.Join(config.BasePath, "lapiasse.log")

	logFileFlags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if env.RunEnv == env.EnvDevelopment {
		logFileFlags |= os.O_TRUNC
	}

	logFile, err := os.OpenFile(logFilePath, logFileFlags, 0666)
	if err != nil {
		return nil, fmt.Errorf("opening log file %q: %w", logFilePath, err)
	}

	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	cleanup(closeLogger(logFile))

	return logger, nil
}

func closeLogger(logFile *os.File) ctxval.CleanupFunc {
	return ctxval.CleanupFunc(func(ctx context.Context) {
		// todo:
		_ = logFile
		applog.Info(ctx, "[TODO] Closing this before other components cause their logs not to be included in the log file before the process stops.")

		// applog.Info(ctx, "Closing log file...")
		// if err := logFile.Close(); err != nil {
		// 	applog.Error(ctx, "Closing log file failed during shutdown", "error", err)
		// } else {
		// applog.Info(ctx, "Closing log file completed")
		// }
	})
}
