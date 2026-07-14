package middleware

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestGlobalErrorHandlerUsesFiberStatusCode(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: GlobalErrorHandler})
	app.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusTeapot, "teapot")
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusTeapot {
		t.Fatalf("expected status %d, got %d", fiber.StatusTeapot, resp.StatusCode)
	}
}

func TestGlobalErrorHandlerDefaultsToInternalServerError(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: GlobalErrorHandler})
	app.Get("/", func(c *fiber.Ctx) error {
		return errors.New("boom")
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(RecoveryMiddleware())
	app.Get("/", func(c *fiber.Ctx) error {
		panic("boom")
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}
}
