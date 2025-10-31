package app

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"adeynack.net/lapiasse/pkg/repository"
	"gorm.io/gorm"
)

// Instance represents a running instance of the application.
type Instance struct {
	DataFileSystem *os.Root // File system for the currently open data "file" (directory).
	DB             *gorm.DB
}

func NewInstance(ch *ConfigurationHolder) (*Instance, error) {
	if ch == nil {
		return nil, errors.New("configuration holder is nil")
	}

	config := ch.Configuration

	slog.Debug("Ensure data directory exists", slog.String("path", config.Data.BasePath))
	if err := os.MkdirAll(config.Data.BasePath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("creating application data directory %q: %w", config.Data.BasePath, err)
	}

	db, err := repository.InitializeGorm(config.Data)
	if err != nil {
		return nil, err
	}

	dataFs, err := os.OpenRoot(config.Data.BasePath)
	if err != nil {
		return nil, fmt.Errorf("opening data file system at %q: %w", config.Data.BasePath, err)
	}

	instance := &Instance{
		DataFileSystem: dataFs,
		DB:             db,
	}

	return instance, nil
}

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

	return errs
}
