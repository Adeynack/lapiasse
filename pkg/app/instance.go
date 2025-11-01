package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"adeynack.net/lapiasse/pkg/repository"
	"gorm.io/gorm"
)

// Instance represents a running instance of the application.
type Instance struct {
	DataFileSystem *os.Root // File system for the currently open data "file" (directory).
	DB             *gorm.DB
	Logger         *slog.Logger

	dependencyContext *ctxval.Resolver
	loggerCloseFn     func() error
}

func NewInstance(ch *ConfigurationHolder) (*Instance, error) {
	if ch == nil {
		return nil, errors.New("configuration holder is nil")
	}

	var err error
	instance := &Instance{
		dependencyContext: ctxval.NewResolver(context.Background()),
	}
	ctx := context.Context(instance.dependencyContext)

	config := ch.Configuration

	instance.Logger, instance.loggerCloseFn, err = configureLogger(config.Data)
	if err != nil {
		return nil, fmt.Errorf("configuring logger: %w", err)
	}
	ctxval.RegisterInResolver(instance.dependencyContext, instance.Logger)

	applog.Debug(ctx, "Ensure data directory exists", "path", config.Data.BasePath)
	if err := os.MkdirAll(config.Data.BasePath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("creating application data directory %q: %w", config.Data.BasePath, err)
	}

	instance.DB, err = repository.InitializeGorm(ctx, config.Data)
	if err != nil {
		return nil, err
	}
	ctxval.RegisterInResolver(instance.dependencyContext, instance.DB)

	instance.DataFileSystem, err = os.OpenRoot(config.Data.BasePath)
	if err != nil {
		return nil, fmt.Errorf("opening data file system at %q: %w", config.Data.BasePath, err)
	}

	return instance, nil
}

// Close implements the [io.Closer] interface.
func (i *Instance) Close() error {
	if i == nil {
		return nil
	}

	// Join errors instead of early-returning the first, in order to attempt closing
	// all resources of the instance.
	var errs error

	if i.DataFileSystem != nil {
		if err := i.DataFileSystem.Close(); err != nil {
			errs = errors.Join(errs, fmt.Errorf("closing data file system: %w", err))
		}
	}

	if i.DB != nil {
		sqlDB, err := i.DB.DB()
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("obtaining sql.DB from gorm.DB: %w", err))
		}

		if err := sqlDB.Close(); err != nil {
			errs = errors.Join(errs, fmt.Errorf("closing database connection: %w", err))
		}
	}

	if i.loggerCloseFn != nil {
		err := i.loggerCloseFn()
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("closing logger: %w", err))
		}
	}

	return errs
}
