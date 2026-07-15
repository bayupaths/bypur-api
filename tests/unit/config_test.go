package unit

import (
	"reflect"
	"testing"

	"github.com/bayupaths/bypur-api/internal/config"
)

func TestParseCorsOriginsTrimsCommaSeparatedValues(t *testing.T) {
	got := config.ParseCorsOrigins("https://bayupur.dev, https://cms.bayupur.dev ,http://localhost:3000")
	want := []string{"https://bayupur.dev", "https://cms.bayupur.dev", "http://localhost:3000"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %#v, got %#v", want, got)
	}
}

func TestEnvironmentHelpers(t *testing.T) {
	cfg := &config.Config{App: config.AppConfig{Env: "production"}}
	if !cfg.IsProduction() || cfg.IsDevelopment() || cfg.IsTest() {
		t.Fatal("production environment helpers returned unexpected values")
	}

	cfg.App.Env = "development"
	if !cfg.IsDevelopment() || cfg.IsProduction() || cfg.IsTest() {
		t.Fatal("development environment helpers returned unexpected values")
	}

	cfg.App.Env = "test"
	if !cfg.IsTest() || cfg.IsProduction() || cfg.IsDevelopment() {
		t.Fatal("test environment helpers returned unexpected values")
	}
}

func TestFeatureFlagsRequireProviderAndFeatureEnabled(t *testing.T) {
	cfg := &config.Config{
		Storage: config.StorageConfig{Enabled: true},
		Mail:    config.MailConfig{Enabled: true},
		Feature: config.FeatureConfig{Storage: false, Email: false},
	}

	if cfg.IsStorageEnabled() {
		t.Fatal("storage should stay disabled until feature flag is enabled")
	}
	if cfg.IsMailEnabled() {
		t.Fatal("mail should stay disabled until feature flag is enabled")
	}

	cfg.Feature.Storage = true
	cfg.Feature.Email = true

	if !cfg.IsStorageEnabled() {
		t.Fatal("storage should be enabled when provider and feature flag are enabled")
	}
	if !cfg.IsMailEnabled() {
		t.Fatal("mail should be enabled when provider and feature flag are enabled")
	}
}
