package repository

import (
	"context"
	"fmt"
	"os"

	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
)

type DataFileSystem struct {
	*os.Root
}

func InitializeDataFilesystem(ctx context.Context, config *Configuration) (*DataFileSystem, error) {
	cleanup, err := ctxval.Resolve[ctxval.CleanupRecorder](ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving cleanup recorder: %w", err)
	}

	applog.Debug(ctx, "Ensure data directory exists", "path", config.BasePath)
	if err := os.MkdirAll(config.BasePath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("creating application data directory %q: %w", config.BasePath, err)
	}

	root, err := os.OpenRoot(config.BasePath)
	if err != nil {
		return nil, fmt.Errorf("opening data file system at %q: %w", config.BasePath, err)
	}
	dataFileSystem := &DataFileSystem{Root: root}

	cleanup(closeDataFilesystem(dataFileSystem))

	return dataFileSystem, nil
}

func closeDataFilesystem(dfs *DataFileSystem) ctxval.CleanupFunc {
	return ctxval.CleanupFunc(func(ctx context.Context) {
		applog.Info(ctx, "Closing data file system...")
		if err := dfs.Close(); err != nil {
			applog.Error(ctx, "Closing data file system failed during shutdown", "error", err)
		} else {
			applog.Info(ctx, "Closing data file system completed")
		}
	})
}
