package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bayupaths/bypur-api/internal/config"
	"github.com/bayupaths/bypur-api/internal/handler"
	"github.com/bayupaths/bypur-api/internal/model"
	"github.com/bayupaths/bypur-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type routerMessageRepo struct {
	created *model.ContactMessage
	err     error
}

func (r *routerMessageRepo) GetMessages(ctx context.Context, status string) ([]model.ContactMessage, error) {
	return nil, r.err
}

func (r *routerMessageRepo) GetByID(ctx context.Context, id string) (*model.ContactMessage, error) {
	return nil, r.err
}

func (r *routerMessageRepo) Create(ctx context.Context, msg *model.ContactMessage) error {
	r.created = msg
	return r.err
}

func (r *routerMessageRepo) Update(ctx context.Context, msg *model.ContactMessage) error {
	return r.err
}

func (r *routerMessageRepo) Delete(ctx context.Context, msg *model.ContactMessage) error {
	return r.err
}

func (r *routerMessageRepo) GetStats(ctx context.Context) (map[string]int64, error) {
	return nil, r.err
}

func newRouterTestApp(cfg *config.Config, msgRepo *routerMessageRepo) *fiber.App {
	app := fiber.New()
	portfolioSvc := service.NewPortfolioService(nil, nil, nil, nil, nil, msgRepo)
	router := &handler.Router{
		App:         app,
		Cfg:         cfg,
		PublicPortH: handler.NewPublicPortfolioHandler(portfolioSvc, nil, cfg),
		AdminH:      &handler.AdminHandler{},
		StorH:       &handler.StorageHandler{},
	}
	router.Setup()
	return app
}

func routerTestConfig(env string) *config.Config {
	return &config.Config{
		App: config.AppConfig{Name: "portfolio-api", Version: "1.2.3", Env: env},
		Server: config.ServerConfig{
			ParsedCorsOrigins: []string{"https://portfolio.example.com", "https://cms.example.com"},
		},
		JWT:      config.JWTConfig{Secret: "01234567890123456789012345678901"},
		Security: config.SecurityConfig{XApiKey: "public-test-key"},
	}
}

func TestRouterHealthAndVersion(t *testing.T) {
	app := newRouterTestApp(routerTestConfig("test"), &routerMessageRepo{})

	healthResp, err := app.Test(httptest.NewRequest(http.MethodGet, "/health", nil))
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	if healthResp.StatusCode != http.StatusOK {
		t.Fatalf("expected health status %d, got %d", http.StatusOK, healthResp.StatusCode)
	}

	versionResp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/version", nil))
	if err != nil {
		t.Fatalf("version request failed: %v", err)
	}
	if versionResp.StatusCode != http.StatusOK {
		t.Fatalf("expected version status %d, got %d", http.StatusOK, versionResp.StatusCode)
	}

	var body struct {
		Data struct {
			App         string `json:"app"`
			Version     string `json:"version"`
			Environment string `json:"environment"`
		} `json:"data"`
	}
	if err := json.NewDecoder(versionResp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode version response: %v", err)
	}
	if body.Data.App != "portfolio-api" || body.Data.Version != "1.2.3" || body.Data.Environment != "test" {
		t.Fatalf("unexpected version payload: %+v", body.Data)
	}
}

func TestRouterPublicRoutesRequireApiKey(t *testing.T) {
	app := newRouterTestApp(routerTestConfig("test"), &routerMessageRepo{})

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/public/profile", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestRouterRejectsInvalidContactPayloadBeforePersistence(t *testing.T) {
	msgRepo := &routerMessageRepo{}
	app := newRouterTestApp(routerTestConfig("test"), msgRepo)

	req := httptest.NewRequest(http.MethodPost, "/api/public/contact", strings.NewReader(`{"name":"B","email":"bad","subject":"Hi","message":"short"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "public-test-key")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
	if msgRepo.created != nil {
		t.Fatal("invalid contact payload should not be persisted")
	}
}

func TestRouterSubmitsValidContactPayload(t *testing.T) {
	msgRepo := &routerMessageRepo{}
	app := newRouterTestApp(routerTestConfig("test"), msgRepo)

	body := bytes.NewBufferString(`{"name":"Bayu Purnomo","email":"bayu@example.com","subject":"Project inquiry","message":"Hello, I want to discuss a production project."}`)
	req := httptest.NewRequest(http.MethodPost, "/api/public/contact", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "public-test-key")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	if msgRepo.created == nil || msgRepo.created.Email != "bayu@example.com" || msgRepo.created.Status != "new" {
		t.Fatalf("valid contact payload was not persisted correctly: %+v", msgRepo.created)
	}
}

func TestRouterAdminRoutesRequireJWT(t *testing.T) {
	app := newRouterTestApp(routerTestConfig("test"), &routerMessageRepo{})

	resp, err := app.Test(httptest.NewRequest(http.MethodPut, "/api/admin/profile", strings.NewReader(`{}`)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestRouterAdminRoutesRejectInvalidJWT(t *testing.T) {
	app := newRouterTestApp(routerTestConfig("test"), &routerMessageRepo{})

	req := httptest.NewRequest(http.MethodPut, "/api/admin/profile", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestRouterProductionRateLimitIsEnabled(t *testing.T) {
	app := newRouterTestApp(routerTestConfig("production"), &routerMessageRepo{})
	var lastStatus int

	for i := 0; i < 101; i++ {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/health", nil))
		if err != nil {
			t.Fatalf("request %d failed: %v", i+1, err)
		}
		lastStatus = resp.StatusCode
	}

	if lastStatus != http.StatusTooManyRequests {
		t.Fatalf("expected production limiter to return %d on request 101, got %d", http.StatusTooManyRequests, lastStatus)
	}
}
