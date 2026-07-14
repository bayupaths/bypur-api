package middleware

import (
	"net/http"

	"bayupur-portofolio-be/internal/config"
	"bayupur-portofolio-be/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// AuthenticateApiKey memverifikasi request publik menggunakan header X-API-KEY
func AuthenticateApiKey(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("x-api-key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if cfg.XApiKey == "" {
			return response.SendError(c, "Server API Key is not configured", "Internal Server Error", http.StatusInternalServerError)
		}

		if apiKey == "" || apiKey != cfg.XApiKey {
			return response.SendError(c, "Invalid or empty API Key", "Unauthorized", http.StatusUnauthorized)
		}

		return c.Next()
	}
}
