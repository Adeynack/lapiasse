package app

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"adeynack.net/lapiasse/pkg/env"
	"adeynack.net/lapiasse/pkg/repository"
)

func configureLogger(config *repository.Configuration) (*slog.Logger, func() error, error) {
	logFilePath := path.Join(config.BasePath, "lapiasse.log")

	logFileFlags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if env.RunEnv == env.EnvDevelopment {
		logFileFlags |= os.O_TRUNC
	}

	logFile, err := os.OpenFile(logFilePath, logFileFlags, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("opening log file %q: %w", logFilePath, err)
	}

	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	closeFn := func() error {
		if err := logFile.Close(); err != nil {
			return fmt.Errorf("closing log file %q: %w", logFilePath, err)
		}

		return nil
	}

	return logger, closeFn, nil
}
