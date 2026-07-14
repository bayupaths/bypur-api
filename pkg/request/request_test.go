package request

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestGetUserId(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		c.Locals("userId", "user-1")

		userID, err := GetUserId(c)
		if err != nil {
			t.Fatalf("GetUserId returned error: %v", err)
		}
		if userID != "user-1" {
			t.Fatalf("expected user-1, got %s", userID)
		}

		return nil
	})

	if _, err := app.Test(httptest.NewRequest("GET", "/", nil)); err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestGetUserIdRequiresAuthenticatedUser(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		if _, err := GetUserId(c); err == nil {
			t.Fatal("expected unauthenticated error")
		}
		return nil
	})

	if _, err := app.Test(httptest.NewRequest("GET", "/", nil)); err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestParsePagination(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		opts := ParsePagination(c, 25)
		if opts.Page != 2 || opts.Limit != 100 || opts.Skip != 100 {
			t.Fatalf("unexpected pagination values: %+v", opts)
		}
		if opts.Search != "api" || opts.SortBy != "name" || opts.SortOrder != "asc" {
			t.Fatalf("unexpected filter values: %+v", opts)
		}
		return nil
	})

	req := httptest.NewRequest("GET", "/?page=2&limit=200&search=%20api%20&sortBy=name&sortOrder=ASC", nil)
	if _, err := app.Test(req); err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestParsePaginationDefaultsInvalidValues(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		opts := ParsePagination(c)
		if opts.Page != 1 || opts.Limit != 10 || opts.Skip != 0 || opts.SortOrder != "desc" {
			t.Fatalf("unexpected default pagination values: %+v", opts)
		}
		return nil
	})

	req := httptest.NewRequest("GET", "/?page=-1&limit=bad&sortOrder=sideways", nil)
	if _, err := app.Test(req); err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

type validationPayload struct {
	Name string `json:"name" validate:"required,min=3"`
}

func TestValidateBody(t *testing.T) {
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		var payload validationPayload
		if err := ValidateBody(c, &payload); err != nil {
			t.Fatalf("ValidateBody returned error: %v", err)
		}
		if payload.Name != "Bayu" {
			t.Fatalf("expected parsed payload name, got %s", payload.Name)
		}
		return nil
	})

	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"Bayu"}`))
	req.Header.Set("Content-Type", "application/json")
	if _, err := app.Test(req); err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestValidateBodyReturnsValidationError(t *testing.T) {
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		var payload validationPayload
		if err := ValidateBody(c, &payload); err == nil {
			t.Fatal("expected validation error")
		}
		return nil
	})

	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"Bo"}`))
	req.Header.Set("Content-Type", "application/json")
	if _, err := app.Test(req); err != nil {
		t.Fatalf("request failed: %v", err)
	}
}
