package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bayupaths/bypur-api/internal/config"
	"github.com/bayupaths/bypur-api/internal/database"
	"github.com/bayupaths/bypur-api/internal/handler"
	"github.com/bayupaths/bypur-api/internal/middleware"
	"github.com/bayupaths/bypur-api/internal/service"
	"github.com/bayupaths/bypur-api/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// Start menjalankan proses bootstrapping aplikasi
func Start() {
	// 1. Muat konfigurasi
	cfg := config.LoadConfig()

	// 2. Inisialisasi Logger Terstruktur
	logger.SetupLogger(cfg.Log.Level, cfg.Log.Format)

	slog.Info("Starting application bootstrap...", "app", cfg.App.Name, "version", cfg.App.Version)

	// 3. Koneksi Database
	db := database.ConnectDB(cfg)

	// 4. Inisialisasi Container (Repositories, Services, Handlers)
	c := initContainer(db, cfg)

	// 5. Verifikasi Integrasi Eksternal secara async
	go verifyExternalIntegrations(c.services.storageService, c.services.mailService)

	// 6. Setup Fiber App
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: middleware.GlobalErrorHandler,
		BodyLimit:    10 * 1024 * 1024, // 10MB limit
	})

	// Setup routing & middlewares
	router := &handler.Router{
		App:         app,
		Cfg:         cfg,
		AuthH:       c.handlers.authH,
		PublicPortH: c.handlers.publicPortH,
		AdminH:      c.handlers.adminH,
		StorH:       c.handlers.storageH,
	}
	router.Setup()

	// 9. Graceful Shutdown & Start Server
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		addr := ":" + strconv.Itoa(cfg.Server.Port)
		slog.Info(fmt.Sprintf("[%s] v%s is running on port %d", cfg.App.Name, cfg.App.Version, cfg.Server.Port))
		slog.Info("Environment: " + cfg.App.Env)

		if err := app.Listen(addr); err != nil && !strings.Contains(err.Error(), "closed") {
			slog.Error("Failed to start Fiber server", "error", err)
			os.Exit(1)
		}
	}()

	sig := <-shutdownChan
	slog.Warn("Shutdown signal received, starting graceful shutdown process...", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("Closing HTTP server...")
	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("Failed to gracefully shut down HTTP server", "error", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		slog.Info("Closing database connection pool...")
		_ = sqlDB.Close()
	}

	slog.Info("Server shutdown completed cleanly. Exiting.")
}

func verifyExternalIntegrations(storage *service.StorageService, mail *service.MailService) {
	slog.Info("Verifying external service integrations...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	isConnected, err := storage.CheckConnection(ctx)
	if err == nil && isConnected {
		slog.Info("Cloudflare R2 Storage integration: CONNECTED")
	} else {
		slog.Warn("R2 Storage connection verification failed - upload features will be affected", "error", err)
	}

	err = mail.VerifyConnection()
	if err == nil {
		slog.Info("SMTP Email integration: CONNECTED")
	} else {
		slog.Warn("SMTP Email connection verification failed - email notifications disabled", "error", err)
	}
}
