package middleware

import (
	"errors"
	"net/http"
	"strings"

	"bayupur-portofolio-be/internal/config"
	"bayupur-portofolio-be/pkg/jwt"
	"bayupur-portofolio-be/pkg/response"

	"github.com/gofiber/fiber/v2"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

// AuthenticateJWT memverifikasi token JWT dari Authorization header (Bearer) atau cookie
func AuthenticateJWT(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := extractToken(c)

		if tokenStr == "" {
			return response.SendError(c, "No access token provided", "Unauthorized", http.StatusUnauthorized)
		}

		claims, err := jwt.VerifyToken(tokenStr, cfg.JWTSecret)
		if err != nil {
			statusText := "Invalid token"
			if errors.Is(err, jwt5.ErrTokenExpired) {
				statusText = "Token expired"
			}
			return response.SendError(c, err.Error(), statusText, http.StatusUnauthorized)
		}

		// Simpan userId dan email di Locals context
		c.Locals("userId", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("userClaims", claims)

		return c.Next()
	}
}

// AuthenticateJWTOptional opsional memverifikasi token JWT, tidak akan gagal jika token tidak ada
func AuthenticateJWTOptional(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := extractToken(c)
		if tokenStr != "" {
			claims, err := jwt.VerifyToken(tokenStr, cfg.JWTSecret)
			if err == nil {
				c.Locals("userId", claims.UserID)
				c.Locals("email", claims.Email)
				c.Locals("userClaims", claims)
			}
		}
		return c.Next()
	}
}

func extractToken(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	cookieVal := c.Cookies("accessToken")
	if cookieVal != "" {
		return cookieVal
	}

	return ""
}
