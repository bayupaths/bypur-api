package config

import (
	"log"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	AppName     string `env:"APP_NAME" envDefault:"bypur-api"`
	AppVersion  string `env:"APP_VERSION" envDefault:"1.0.0"`
	Env         string `env:"ENV" envDefault:"development"` // development, production, test
	Port        int    `env:"PORT" envDefault:"3001"`
	ApiUrl      string `env:"API_URL" envDefault:"http://localhost:3001"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	DirectURL   string `env:"DIRECT_URL"`

	JWTSecret         string `env:"JWT_SECRET,required"`
	JWTAccessExpire   string `env:"JWT_ACCESS_EXPIRE" envDefault:"15m"`
	JWTRefreshExpire  string `env:"JWT_REFRESH_EXPIRE" envDefault:"7d"`
	CorsOrigins       string `env:"CORS_ORIGIN" envDefault:"http://localhost:3000"`
	ParsedCorsOrigins []string

	R2Enabled         bool   `env:"R2_ENABLED" envDefault:"false"`
	R2AccountID       string `env:"R2_ACCOUNT_ID"`
	R2AccessKeyID     string `env:"R2_ACCESS_KEY_ID"`
	R2AccessKeySecret string `env:"R2_ACCESS_KEY_SECRET"`
	R2BucketName      string `env:"R2_BUCKET_NAME"`
	R2PublicURL       string `env:"R2_PUBLIC_URL"`

	SMTPEnabled bool   `env:"SMTP_ENABLED" envDefault:"false"`
	SMTPHost    string `env:"SMTP_HOST" envDefault:"smtp.gmail.com"`
	SMTPPort    int    `env:"SMTP_PORT" envDefault:"587"`
	SMTPUser    string `env:"SMTP_USER"`
	SMTPPass    string `env:"SMTP_PASS"`
	SMTPFrom    string `env:"SMTP_FROM" envDefault:"noreply@bayu-apps.com"`

	LogLevel  string `env:"LOG_LEVEL" envDefault:"debug"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"json"`

	PortfolioURL string `env:"PORTFOLIO_URL" envDefault:"http://localhost:3000"`
	CmsURL       string `env:"CMS_URL" envDefault:"http://localhost:3002"`

	FeatureStorage bool `env:"FEATURE_STORAGE" envDefault:"false"`
	FeatureEmail   bool `env:"FEATURE_EMAIL" envDefault:"false"`

	XApiKey string `env:"X_API_KEY" envDefault:"bypur-default-public-api-key-2026"`
}

func LoadConfig() *Config {
	// Membaca file .env jika ada
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to load ENV configuration: %v", err)
	}

	if len(cfg.JWTSecret) < 32 {
		log.Fatalf("JWT_SECRET is too short! It must be at least 32 characters long for security.")
	}

	// Parsing CORS Origins (comma-separated list)
	origins := strings.Split(cfg.CorsOrigins, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}
	cfg.ParsedCorsOrigins = origins

	return cfg
}
