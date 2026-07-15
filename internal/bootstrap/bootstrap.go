package bootstrap

import (
	"log/slog"

	"github.com/bayupaths/bypur-api/internal/config"
	"github.com/bayupaths/bypur-api/internal/database"
	"github.com/bayupaths/bypur-api/pkg/logger"

	"gorm.io/gorm"
)

func Start() {
	cfg := config.LoadConfig()
	logger.SetupLogger(cfg.Log.Level, cfg.Log.Format, cfg.Log.Dir)
	slog.Info("Starting application bootstrap...", "app", cfg.App.Name, "version", cfg.App.Version)

	db := setupDatabase(cfg)
	c := initContainer(db, cfg)

	go verifyExternalIntegrations(c.services.storageService, c.services.mailService)

	runServer(cfg, c)

	database.Close(db)
	logger.Close()
	slog.Info("Server shutdown completed cleanly. Exiting.")
}

func setupDatabase(cfg *config.Config) *gorm.DB {
	db := database.ConnectDB(cfg)
	go database.StartKeepAlive(db)
	return db
}
