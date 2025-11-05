package repository

import (
	"context"
	"fmt"
	"log/slog"

	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/model"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeGorm(ctx context.Context, config *Configuration) (*gorm.DB, error) {
	applog.Debug(ctx, "Initializing the Gorm database connector", "main_database_file_path", config.MainDatabaseFilePath())

	gormLogger := slogGorm.New(
		slogGorm.WithTraceAll(),
		slogGorm.WithContextFunc("system", func(context.Context) (slog.Value, bool) {
			return slog.StringValue("gorm"), true
		}),
	)

	gormConfig := &gorm.Config{
		Logger: gormLogger,
	}

	db, err := gorm.Open(sqlite.Open(config.MainDatabaseFilePath()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("opening main database %q: %w", config.MainDatabaseFilePath(), err)
	}

	applog.Info(ctx, "Auto-migrating the Gorm models")
	if err := db.AutoMigrate(model.Models...); err != nil {
		return nil, fmt.Errorf("migrating database schema: %w", err)
	}

	ctxval.MustCleanup(ctx, closeDB(db))

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
