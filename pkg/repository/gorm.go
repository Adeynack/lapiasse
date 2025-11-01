package repository

import (
	"context"
	"fmt"
	"log/slog"

	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/model"
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

	applog.Debug(ctx, "Auto-migrating the Gorm models")
	if err := db.AutoMigrate(model.Models...); err != nil {
		return nil, fmt.Errorf("migrating database schema: %w", err)
	}

	return db, nil
}
