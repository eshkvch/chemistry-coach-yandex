package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dsn string, isDev bool) (*gorm.DB, error) {
	logLevel := logger.Silent
	if isDev {
		logLevel = logger.Warn
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logLevel)})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	return db, nil
}

func RunMigrations(db *gorm.DB, migrationsDir string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	path := filepath.Join(migrationsDir, "001_init.sql")
	body, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}
	if _, err := sqlDB.Exec(string(body)); err != nil {
		return fmt.Errorf("exec migration: %w", err)
	}
	return nil
}

func Ping(db *gorm.DB) error {
	var sqlDB *sql.DB
	var err error
	sqlDB, err = db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
