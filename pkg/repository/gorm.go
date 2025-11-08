package repository

import (
	"context"
	"fmt"
	"log/slog"

	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/model"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"github.com/go-chi/chi/v5/middleware"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeGorm(ctx context.Context, config *Configuration) (*gorm.DB, error) {
	dsn := config.MainDatabaseFilePath()
	applog.Debug(ctx, "Initializing the Gorm database connector", "main_database_file_path", dsn)

	gormLogger, err := initializeLogger(ctx)
	if err != nil {
		return nil, fmt.Errorf("initializing Gorm logger: %w", err)
	}

	gormConfig := &gorm.Config{
		Logger: gormLogger,
	}

	db, err := gorm.Open(sqlite.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("opening main database %q: %w", dsn, err)
	}

	ctxval.MustCleanup(ctx, closeDB(db))

	applog.Info(ctx, "Auto-migrating the Gorm models")
	if err := db.AutoMigrate(model.Models...); err != nil {
		return nil, fmt.Errorf("migrating database schema: %w", err)
	}

	return db, nil
}

func closeDB(db *gorm.DB) ctxval.CleanupFunc {
	return ctxval.CleanupFunc(func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("getting underlying SQL DB from Gorm DB failed during shutdown: %w", err)
		}

		applog.Info(ctx, "Closing Gorm database connection...")
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("closing Gorm database connection failed during shutdown: %w", err)
		}

		return nil
	})
}

func initializeLogger(ctx context.Context) (logger.Interface, error) {
	logger, err := applog.FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting application logger from context: %w", err)
	}

	gormLogger := slogGorm.New(
		slogGorm.WithHandler(logger.Handler()),
		slogGorm.WithTraceAll(),
		withMetaGroupName(),
		slogGorm.WithContextValue("request_id", middleware.RequestIDKey),
	)

	return gormLogger, nil
}

func withMetaGroupName() slogGorm.Option {
	return slogGorm.WithContextFunc("_group", func(context.Context) (slog.Value, bool) {
		return slog.StringValue("gorm"), true
	})
}
