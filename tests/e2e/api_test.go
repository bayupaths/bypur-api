package e2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bayupaths/bypur-api/internal/config"
	"github.com/bayupaths/bypur-api/internal/handler"
	"github.com/bayupaths/bypur-api/internal/model"
	"github.com/bayupaths/bypur-api/internal/repository"
	"github.com/bayupaths/bypur-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

func TestAPI_GetProfile(t *testing.T) {
	// Initialize Fiber app
	app := fiber.New()

	// Initialize mock repositories/services
	mockSkillRepo := &mockSkillRepository{}
	mockProfileRepo := &mockProfileRepository{
		profile: &model.Profile{
			Name:  "Bayu Purnomo",
			Email: "bayu@example.com",
			Title: "Fullstack Developer",
		},
	}

	portfolioService := service.NewPortfolioService(mockProfileRepo, nil, mockSkillRepo, nil, nil, nil)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			XApiKey: "test-api-key",
		},
		Server: config.ServerConfig{
			ParsedCorsOrigins: []string{"http://localhost:3000"},
		},
	}

	// Setup Router
	router := &handler.Router{
		App:         app,
		Cfg:         cfg,
		PublicPortH: handler.NewPublicPortfolioHandler(portfolioService, nil, cfg),
	}
	router.Setup()

	// Perform HTTP GET request to /api/public/profile
	req := httptest.NewRequest("GET", "/api/public/profile", nil)
	req.Header.Set("x-api-key", "test-api-key")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// Mock repositories needed for API test
type mockProfileRepository struct {
	repository.ProfileRepository
	profile *model.Profile
}

func (m *mockProfileRepository) GetProfile(ctx context.Context) (*model.Profile, error) {
	return m.profile, nil
}

type mockSkillRepository struct {
	repository.SkillRepository
}

func (m *mockSkillRepository) GetSkills(ctx context.Context, category string) ([]model.Skill, error) {
	return []model.Skill{}, nil
}
