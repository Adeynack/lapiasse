package repository

import (
	"fmt"

	"adeynack.net/lapiasse/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeGorm(config *Configuration) (*gorm.DB, error) {
	gormConfig := &gorm.Config{}

	db, err := gorm.Open(sqlite.Open(config.MainDatabaseFilePath()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("opening main database %q: %w", config.MainDatabaseFilePath(), err)
	}

	if err := db.AutoMigrate(model.Models...); err != nil {
		return nil, fmt.Errorf("migrating database schema: %w", err)
	}

	return db, nil
}
