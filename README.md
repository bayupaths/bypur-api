# Bayupur Portfolio API (Back-End)

This is the Go-based backend REST API for developer portfolios and CMS applications, built using Fiber, GORM, PostgreSQL, Cloudflare R2 object storage, and golang-migrate.

## Technology Stack

- **Core**: Go (Golang)
- **Web Framework**: [Fiber v2](https://github.com/gofiber/fiber)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL (via GORM/pgx driver)
- **Database Migrations**: [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- **Storage Integration**: Cloudflare R2 Storage (AWS S3-compatible SDK v2)
- **Validation**: [validator/v10](https://github.com/go-playground/validator)

---

## Key Architectural Decisions

1. **Repository Pattern Layer**:
   All direct GORM/database operations are decoupled from services and placed under `internal/repository`. Services reference interface types, allowing clean mocking during unit tests.

2. **Modular Handlers (SRP)**:
   CMS/admin controllers are split into domain-specific modules (`admin_profile`, `admin_experience`, `admin_offering`, `admin_skill`, `admin_project`, `admin_contact`) to enforce Single Responsibility Principle and maintainability.

3. **Versioned Database Migrations**:
   We use structured SQL migration scripts under `internal/database/migrations` instead of GORM's `AutoMigrate`. These migrations are automatically packaged inside the compiled binary using Go `embed` and executed during startup.

---

## Local Setup & Development

### 1. Prerequisites
- **Go** (version 1.25 or later recommended)
- **PostgreSQL** instance or Docker running local DB
- **Cloudflare R2** account (optional for local mock environment, or disable in env configurations)

### 2. Configuration Setup
Copy the environment template and set your postgres credentials:
```bash
cp .env.example .env
```

Ensure `DATABASE_URL` is mapped to your PostgreSQL instance:
```env
DATABASE_URL=postgresql://postgres:[password]@localhost:5432/postgres?sslmode=disable
```

### 3. Spin Up Local Services (Optional)
If you want to use the bundled PostgreSQL via Docker Compose:
```bash
docker-compose up -d
```

### 4. Running the Application
Run the main server binary. Migrations will run automatically:
```bash
go run cmd/api/main.go
```

The server runs on `http://localhost:3001` by default.

---

## Directory Layout

- `cmd/api/main.go`: Application entrypoint.
- `internal/`:
  - `bootstrap/`: Application lifecycle hooks, dependency configuration, and integrations checklist.
  - `config/`: Configuration mapping via env variables.
  - `database/`: Database pool initialization and embedded migration runner.
  - `handler/`: Public and modular admin Fiber HTTP handler controllers.
  - `middleware/`: Rate limiters, JWT authorization checks, recovery, and security middleware.
  - `model/`: GORM database structs.
  - `repository/`: Database query and transaction operations.
  - `service/`: Domain-specific business logic (JWT auth, Cloudflare storage, SMTP email).
- `pkg/`: Shared utility packages (JWT helpers, request parse validators, standard JSON responses).
