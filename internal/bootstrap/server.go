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
	"github.com/bayupaths/bypur-api/internal/handler"
	"github.com/bayupaths/bypur-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func runServer(cfg *config.Config, c *container) {
	app := newFiberApp(cfg)
	setupRoutes(app, cfg, c)
	startListening(app, cfg)
	waitForShutdown(app)
}

func newFiberApp(cfg *config.Config) *fiber.App {
	return fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: middleware.GlobalErrorHandler,
		BodyLimit:    10 * 1024 * 1024,
	})
}

func setupRoutes(app *fiber.App, cfg *config.Config, c *container) {
	router := &handler.Router{
		App:         app,
		Cfg:         cfg,
		AuthH:       c.handlers.authH,
		PublicPortH: c.handlers.publicPortH,
		AdminH:      c.handlers.adminH,
		StorH:       c.handlers.storageH,
	}
	router.Setup()
}

func startListening(app *fiber.App, cfg *config.Config) {
	go func() {
		addr := ":" + strconv.Itoa(cfg.Server.Port)
		slog.Info(
			fmt.Sprintf("[%s] v%s is running", cfg.App.Name, cfg.App.Version),
			"port", cfg.Server.Port,
			"env", cfg.App.Env,
		)
		if err := app.Listen(addr); err != nil && !strings.Contains(err.Error(), "closed") {
			slog.Error("Failed to start Fiber server", "error", err)
			os.Exit(1)
		}
	}()
}

func waitForShutdown(app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	sig := <-quit
	slog.Warn("Shutdown signal received, starting graceful shutdown...", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("Closing HTTP server...")
	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("Failed to gracefully shut down HTTP server", "error", err)
	}
}
