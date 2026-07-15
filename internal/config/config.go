package config

import (
	"log"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Name    string `env:"NAME" envDefault:"bypur-api"`
	Version string `env:"VERSION" envDefault:"1.0.0"`
	Env     string `env:"ENV" envDefault:"development"`
}

type ServerConfig struct {
	Port              int    `env:"PORT" envDefault:"3001"`
	ApiUrl            string `env:"API_URL" envDefault:"http://localhost:3001"`
	CorsOrigins       string `env:"CORS_ORIGIN" envDefault:"http://localhost:3000"`
	ParsedCorsOrigins []string
}

type DBConfig struct {
	URL       string `env:"URL,required"`
	DirectURL string `env:"DIRECT_URL"`
}

type JWTConfig struct {
	Secret        string `env:"SECRET,required"`
	AccessExpire  string `env:"ACCESS_EXPIRE" envDefault:"15m"`
	RefreshExpire string `env:"REFRESH_EXPIRE" envDefault:"7d"`
}

type StorageConfig struct {
	Enabled   bool   `env:"ENABLED" envDefault:"false"`
	AccountID string `env:"ACCOUNT_ID"`
	AccessKey string `env:"ACCESS_KEY_ID"`
	SecretKey string `env:"ACCESS_KEY_SECRET"`
	Bucket    string `env:"BUCKET_NAME"`
	PublicURL string `env:"PUBLIC_URL"`
}

type MailConfig struct {
	Enabled bool   `env:"ENABLED" envDefault:"false"`
	Host    string `env:"HOST" envDefault:"smtp.gmail.com"`
	Port    int    `env:"PORT" envDefault:"587"`
	User    string `env:"USER"`
	Pass    string `env:"PASS"`
	From    string `env:"FROM" envDefault:"noreply@bayu-apps.com"`
}

type LogConfig struct {
	Level  string `env:"LEVEL" envDefault:"debug"`
	Format string `env:"FORMAT" envDefault:"json"`
}

type FrontendConfig struct {
	PortfolioURL string `env:"PORTFOLIO_URL" envDefault:"http://localhost:3000"`
	CmsURL       string `env:"CMS_URL" envDefault:"http://localhost:3002"`
}

type FeatureConfig struct {
	Storage bool `env:"STORAGE" envDefault:"false"`
	Email   bool `env:"EMAIL" envDefault:"false"`
}

type SecurityConfig struct {
	XApiKey string `env:"X_API_KEY" envDefault:"bypur-default-public-api-key-2026"`
}

type Config struct {
	App      AppConfig      `envPrefix:"APP_"`
	Server   ServerConfig   `envPrefix:"SERVER_"`
	DB       DBConfig       `envPrefix:"DB_"`
	JWT      JWTConfig      `envPrefix:"JWT_"`
	Storage  StorageConfig  `envPrefix:"STORAGE_"`
	Mail     MailConfig     `envPrefix:"MAIL_"`
	Log      LogConfig      `envPrefix:"LOG_"`
	Frontend FrontendConfig `envPrefix:"FRONTEND_"`
	Feature  FeatureConfig  `envPrefix:"FEATURE_"`
	Security SecurityConfig `envPrefix:"SECURITY_"`
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to load ENV configuration: %v", err)
	}

	if len(cfg.JWT.Secret) < 32 {
		log.Fatalf("JWT_SECRET is too short! It must be at least 32 characters long for security.")
	}

	cfg.Server.ParsedCorsOrigins = ParseCorsOrigins(cfg.Server.CorsOrigins)

	return cfg
}

func ParseCorsOrigins(raw string) []string {
	origins := strings.Split(raw, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}
	return origins
}

func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

func (c *Config) IsTest() bool {
	return c.App.Env == "test"
}

func (c *Config) IsStorageEnabled() bool {
	return c.Storage.Enabled && c.Feature.Storage
}

func (c *Config) IsMailEnabled() bool {
	return c.Mail.Enabled && c.Feature.Email
}
