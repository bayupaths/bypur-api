package bootstrap

import (
	"context"
	"log/slog"
	"time"

	"github.com/bayupaths/bypur-api/internal/service"
)

func verifyExternalIntegrations(storage *service.StorageService, mail *service.MailService) {
	slog.Info("Verifying external service integrations...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	verifyStorage(ctx, storage)
	verifyMail(mail)
}

func verifyStorage(ctx context.Context, storage *service.StorageService) {
	isConnected, err := storage.CheckConnection(ctx)
	if err == nil && isConnected {
		slog.Info("Cloudflare R2 Storage integration: CONNECTED")
	} else {
		slog.Warn("R2 Storage connection verification failed — upload features will be affected", "error", err)
	}
}

func verifyMail(mail *service.MailService) {
	if err := mail.VerifyConnection(); err != nil {
		slog.Warn("SMTP Email connection verification failed — email notifications disabled", "error", err)
	} else {
		slog.Info("SMTP Email integration: CONNECTED")
	}
}
