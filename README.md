<div align="center">

# 🚀 Bypur Portfolio API

**Production-ready REST API for portfolio showcase & CMS management**

Built with Go, Fiber v2, GORM, PostgreSQL, and Cloudflare R2

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)
![Fiber](https://img.shields.io/badge/Fiber-v2-00ACD7?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supabase-336791?logo=postgresql)
![License](https://img.shields.io/badge/License-MIT-yellow)

[Features](#-features) • [Tech Stack](#-tech-stack) • [Quick Start](#-quick-start) • [Project Structure](#-project-structure) • [Environment](#️-environment-configuration)

</div>

---

## ✨ Features

- 🔐 **JWT Authentication** — Access & refresh token with secure HttpOnly cookie handling
- 🛡️ **Security First** — API key guard, CORS, account lockout after failed login attempts
- 📂 **Cloudflare R2 Storage** — S3-compatible object storage with presigned URLs
- 📧 **SMTP Email** — Verification, password reset & welcome email via styled HTML templates
- 🗄️ **Auto Migration** — SQL migrations embedded in binary via `go:embed`, run on startup
- 📊 **Structured Logging** — `log/slog` with file rotation (lumberjack) to `logs/app.log`
- 🩺 **DB Keep-Alive** — Periodic ping to prevent Supabase free-tier from pausing
- 🧪 **Test Pyramid** — Unit, integration & E2E tests with SonarQube coverage report
- 🐳 **Containerized** — Docker & Docker Compose ready
- 🔄 **CI/CD** — Jenkins pipeline & GitHub Actions (CI + Supabase keep-alive)

## 🛠 Tech Stack

```yaml
Runtime:    Go 1.25+
Framework:  Fiber v2
ORM:        GORM (pgx driver)
Database:   PostgreSQL (Supabase)
Migrations: golang-migrate/migrate (embedded SQL)
Storage:    Cloudflare R2 (AWS S3-compatible SDK v2)
Email:      SMTP (net/smtp)
Auth:       JWT (golang-jwt/jwt)
Logging:    log/slog + lumberjack (file rotation)
Validation: go-playground/validator v10
Testing:    Go testing + testify
CI/CD:      Jenkins + GitHub Actions
Container:  Docker + Docker Compose
```

## 🚀 Quick Start

### Prerequisites

- Go **1.25+**
- PostgreSQL instance — local, Docker, or [Supabase](https://supabase.com)
- Cloudflare R2 account *(optional — disable with `FEATURE_STORAGE=false`)*

### Installation

```bash
# Clone repository
git clone https://github.com/bayupaths/bypur-api.git
cd bypur-api

# Setup environment variables
cp .env.example .env
# Edit .env with your configuration

# Run the application (migrations run automatically on startup)
go run cmd/api/main.go
```

🎉 Server runs on **http://localhost:3001**

### Using Docker

```bash
# Build and run via Docker Compose
docker compose up -d

# View logs
docker compose logs -f app
```

## 📁 Project Structure

```
.
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── bootstrap/
│   │   ├── bootstrap.go         # Orchestrator — calls setup functions in order
│   │   ├── container.go         # DI container (repos → services → handlers)
│   │   ├── server.go            # Fiber HTTP lifecycle (setup, listen, shutdown)
│   │   └── integrations.go      # External service health checks (R2, SMTP)
│   ├── config/
│   │   └── config.go            # Nested env config structs with prefix mapping
│   ├── database/
│   │   ├── database.go          # ConnectDB, Close, StartKeepAlive
│   │   └── migrations/          # Versioned SQL migration files (embedded)
│   ├── dto/
│   │   └── dto.go               # Request/response Data Transfer Objects
│   ├── handler/
│   │   ├── router.go            # Route registration & middleware setup
│   │   ├── auth.go              # Auth endpoints (login, register, refresh, logout)
│   │   ├── public.go            # Public portfolio endpoints (no auth required)
│   │   ├── admin.go             # Admin CMS endpoints (JWT protected)
│   │   └── storage.go           # File upload endpoints
│   ├── middleware/
│   │   ├── auth.go              # JWT verification middleware
│   │   ├── api_key.go           # X-API-Key guard for public routes
│   │   └── error.go             # Global Fiber error handler
│   ├── model/                   # GORM model structs
│   ├── repository/              # Database query layer (interface + implementation)
│   └── service/
│       ├── auth.go              # Auth business logic (JWT, bcrypt, lockout)
│       ├── mail.go              # SMTP email with styled HTML templates
│       ├── portfolio.go         # Portfolio data service
│       └── storage.go           # Cloudflare R2 operations
├── pkg/
│   ├── jwt/                     # JWT helper (sign, verify, parse)
│   ├── logger/                  # Custom slog logger with lumberjack rotation
│   ├── request/                 # Request body parsing & validation helpers
│   └── response/                # Standardized JSON response builders
├── tests/
│   ├── unit/                    # Unit tests (mocked dependencies)
│   ├── integration/             # Integration tests (real DB)
│   └── e2e/                     # End-to-end API tests
├── .github/workflows/
│   ├── ci.yml                   # Build & test on push/PR
│   └── supabase-keepalive.yml   # Scheduled ping every 5 days
├── Dockerfile
├── docker-compose.yml
├── Jenkinsfile
└── sonar-project.properties
```

## 🗄️ Database

**9 Core Models** managed via versioned SQL migrations:

```
User, RefreshToken, Profile, SocialLink,
Offering, Skill, Experience, Project, ContactMessage
```

Migrations are embedded in the binary using `//go:embed` and run automatically on startup via `golang-migrate`. No manual migration commands needed.

### Key Architecture Decisions

1. **Repository Pattern** — All GORM operations are behind interfaces in `internal/repository`, enabling clean mock injection for unit tests.
2. **Embedded Migrations** — SQL files in `internal/database/migrations/` are compiled into the binary — no external migration tooling needed at runtime.
3. **Nested Config Structs** — Config is organized as `cfg.DB.URL`, `cfg.JWT.Secret`, `cfg.Mail.Host` etc. with env prefix mapping (e.g. `DB_URL`, `JWT_SECRET`).
4. **Single Responsibility Bootstrap** — `Start()` is a thin orchestrator; each concern (server, DB, integrations) lives in its own file.

## 🔧 Available Commands

### Development

```bash
go run cmd/api/main.go       # Start development server
go build -o api cmd/api/main.go  # Build production binary
./api                        # Run production binary
```

### Testing

```bash
go test ./...                         # Run all tests
go test ./tests/unit/...              # Unit tests only
go test ./tests/integration/...       # Integration tests only
go test ./tests/e2e/...               # E2E tests only
go test -coverprofile=coverage.out ./...  # With coverage
```

### Docker

```bash
docker compose up -d         # Start application
docker compose down          # Stop all containers
docker compose logs -f app   # Tail application logs
```

## ⚙️ Environment Configuration

Copy and customize the template:

```bash
cp .env.example .env
```

Key variables:

```env
# Application
APP_NAME=bypur-api
APP_ENV=development              # development | production | test

# Server
SERVER_PORT=3001
SERVER_CORS_ORIGIN=http://localhost:3000,http://localhost:3002

# Database (Supabase / PostgreSQL)
DB_URL=postgresql://postgres:[password]@localhost:5432/postgres?sslmode=disable

# JWT Authentication
JWT_SECRET=your_secure_jwt_secret_key_32_bytes_min
JWT_ACCESS_EXPIRE=15m
JWT_REFRESH_EXPIRE=7d

# Cloudflare R2 Storage
STORAGE_ENABLED=false
STORAGE_ACCOUNT_ID=your_cloudflare_account_id
STORAGE_BUCKET_NAME=your_bucket_name
STORAGE_PUBLIC_URL=https://...

# Email (SMTP)
MAIL_ENABLED=false
MAIL_HOST=smtp.gmail.com
MAIL_USER=your_email@gmail.com
MAIL_PASS=your_app_password

# Logging
LOG_LEVEL=debug                  # debug | info | warn | error
LOG_FORMAT=json                  # json | text
LOG_DIR=logs                     # Output directory for log files

# Feature Flags
FEATURE_STORAGE=false
FEATURE_EMAIL=false

# Security
SECURITY_X_API_KEY=bypur-default-public-api-key-2026
```

> 💡 **Tip:** Copy from `.env.example` for the complete variable list with descriptions.

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes using conventional commits:
   - `feat:` New feature
   - `fix:` Bug fix
   - `refactor:` Code restructuring
   - `docs:` Documentation updates
   - `test:` Test additions or updates
   - `ci:` CI/CD changes
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

---

<div align="center">

**Built with ❤️ by [Bayu Purnomo](https://github.com/bayupaths)**

⭐ Star this repo if you find it helpful!

</div>
