package database

import (
	"embed"
	"log/slog"
	"time"

	"github.com/bayupaths/bypur-api/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

//go:embed migrations/*.sql
var migrationFS embed.FS

func ConnectDB(cfg *config.Config) *gorm.DB {
	var logMode logger.Interface
	if cfg.IsDevelopment() {
		logMode = logger.Default.LogMode(logger.Info)
	} else {
		logMode = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(cfg.DB.URL), &gorm.Config{
		Logger: logMode,
	})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to retrieve SQL DB instance", "error", err)
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	slog.Info("PostgreSQL database connection successfully established")

	runMigrations(cfg.DB.URL)

	DB = db
	return db
}

func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		slog.Warn("Failed to retrieve SQL DB instance for closing", "error", err)
		return
	}
	slog.Info("Closing database connection pool...")
	if err := sqlDB.Close(); err != nil {
		slog.Warn("Failed to close database connection pool", "error", err)
	}
}

func StartKeepAlive(db *gorm.DB) {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if err := db.Exec("SELECT 1").Error; err != nil {
			slog.Warn("DB keep-alive ping failed", "error", err)
		} else {
			slog.Info("DB keep-alive ping successful")
		}
	}
}

func runMigrations(databaseURL string) {
	slog.Info("Running database migrations...")
	d, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		slog.Error("Failed to initialize migration iofs driver", "error", err)
		panic(err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		slog.Error("Failed to initialize migrate instance", "error", err)
		panic(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Failed to run database migrations", "error", err)
		panic(err)
	}

	slog.Info("Database migrations completed successfully")
}
