package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bayupaths/bypur-api/internal/config"
	tokenjwt "github.com/bayupaths/bypur-api/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateApiKey(t *testing.T) {
	app := fiber.New()
	app.Get("/", AuthenticateApiKey(&config.Config{Security: config.SecurityConfig{XApiKey: "secret"}}), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("x-api-key", "secret")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected status %d, got %d", fiber.StatusNoContent, resp.StatusCode)
	}
}

func TestAuthenticateApiKeyRejectsInvalidKey(t *testing.T) {
	app := fiber.New()
	app.Get("/", AuthenticateApiKey(&config.Config{Security: config.SecurityConfig{XApiKey: "secret"}}), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAuthenticateJWT(t *testing.T) {
	token, err := tokenjwt.GenerateToken("user-1", "bayu@example.com", "secret", time.Hour, "portfolio")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	app := fiber.New()
	app.Get("/", AuthenticateJWT(&config.Config{JWT: config.JWTConfig{Secret: "secret"}}), func(c *fiber.Ctx) error {
		if c.Locals("userId") != "user-1" || c.Locals("email") != "bayu@example.com" {
			t.Fatalf("JWT locals were not set")
		}
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected status %d, got %d", fiber.StatusNoContent, resp.StatusCode)
	}
}

func TestAuthenticateJWTOptionalIgnoresMissingToken(t *testing.T) {
	app := fiber.New()
	app.Get("/", AuthenticateJWTOptional(&config.Config{JWT: config.JWTConfig{Secret: "secret"}}), func(c *fiber.Ctx) error {
		if c.Locals("userId") != nil {
			t.Fatalf("expected no user local for missing optional token")
		}
		return c.SendStatus(fiber.StatusNoContent)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected status %d, got %d", fiber.StatusNoContent, resp.StatusCode)
	}
}

func TestExtractTokenFromCookie(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		if token := extractToken(c); token != "cookie-token" {
			t.Fatalf("expected cookie token, got %s", token)
		}
		return nil
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "accessToken", Value: "cookie-token"})
	if _, err := app.Test(req); err != nil {
		t.Fatalf("request failed: %v", err)
	}
}
