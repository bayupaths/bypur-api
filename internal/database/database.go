package database

import (
	"embed"
	"log/slog"
	"time"

	"bayupur-portofolio-be/internal/config"

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
	if cfg.Env == "development" {
		logMode = logger.Default.LogMode(logger.Info)
	} else {
		logMode = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logMode,
	})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		panic(err)
	}

	// Mengatur connection pool
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to retrieve SQL DB instance", "error", err)
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	slog.Info("PostgreSQL database connection successfully established")

	// Run golang-migrate database migrations
	runMigrations(cfg.DatabaseURL)

	DB = db
	return db
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
