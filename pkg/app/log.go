package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/env"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"adeynack.net/lapiasse/pkg/platform/slogex"
	"adeynack.net/lapiasse/pkg/repository"
	"github.com/golang-cz/devslog"
)

type logConfiguration struct {
	Verbose       bool `json:"verbose"`         // Switch log level from "info" to "debug".
	ForceJsonFile bool `json:"force_json_file"` // Force log file to be in JSON format even when running in UI-less mode.

	UILess bool `json:"-"` // Application runs without UI (eg: server mode). Log also to STDERR for direct feedback. Non-configurable by config file.
}

func configureLogger(ctx context.Context, config *logConfiguration) (*slog.Logger, error) {
	handlerCandidates := []func(context.Context, *logConfiguration) (slog.Handler, error){
		initializeJsonFileHandler,
		initializeJsonStdOutHandler,
		initializePrettyStdOutHandler,
	}

	handlers := make([]slog.Handler, 0, len(handlerCandidates))
	for _, initHandler := range handlerCandidates {
		h, err := initHandler(ctx, config)
		if err != nil {
			return nil, err
		}
		if h != nil {
			handlers = append(handlers, h)
		}
	}

	var handler slog.Handler
	if len(handlers) == 1 {
		handler = handlers[0]
	} else {
		handler = slogex.NewMultiplexHandler(handlers...)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger, nil
}

func levelFromConfig(config *logConfiguration) slog.Level {
	if config.Verbose {
		return slog.LevelDebug
	}

	return slog.LevelInfo
}

func initializeJsonFileHandler(ctx context.Context, config *logConfiguration) (slog.Handler, error) {
	if config.UILess && !config.ForceJsonFile {
		return nil, nil // in UI-less mode, using only stdout handlers, unless forced to use JSON file by configuration
	}

	dataRoot, err := ctxval.Resolve[*repository.DataFileSystem](ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving data root for logger configuration: %w", err)
	}

	logFileFlags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if env.GetRunEnv(ctx) == env.EnvDevelopment {
		logFileFlags |= os.O_TRUNC
	}

	logFile, err := dataRoot.OpenFile("lapiasse.log", logFileFlags, 0666)
	if err != nil {
		return nil, fmt.Errorf("opening log file: %w", err)
	}

	ctxval.MustCleanup(ctx, ctxval.CleanupFunc(func(ctx context.Context) error {
		applog.Info(ctx, "Closing log file. No further log entries will be written.")
		if err := logFile.Close(); err != nil {
			return fmt.Errorf("closing log file: %w", err)
		}

		return nil
	}))

	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: levelFromConfig(config),
	})

	return handler, nil
}

func initializeJsonStdOutHandler(ctx context.Context, config *logConfiguration) (slog.Handler, error) {
	if !config.UILess {
		return nil, nil // in UI mode, using only JSON file output
	}

	if env.GetRunEnv(ctx) == env.EnvDevelopment {
		return nil, nil // in development, using pretty output
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: levelFromConfig(config),
	})

	return handler, nil
}

func initializePrettyStdOutHandler(ctx context.Context, config *logConfiguration) (slog.Handler, error) {
	if !config.UILess {
		return nil, nil // in UI mode, using only JSON file output
	}

	if env.GetRunEnv(ctx) != env.EnvDevelopment {
		return nil, nil // using pretty output only in development
	}

	handler := devslog.NewHandler(os.Stdout, &devslog.Options{
		HandlerOptions: &slog.HandlerOptions{
			Level: levelFromConfig(config),
		},
		NewLineAfterLog: true,
	})

	return handler, nil
}
