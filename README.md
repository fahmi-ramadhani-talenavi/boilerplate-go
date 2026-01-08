# Go Banking-Grade Boilerplate

A production-ready Go boilerplate following Clean Architecture principles with banking-grade security and observability.

## Features

### Security
- **Security Headers**: HSTS, CSP, X-Frame-Options, X-Content-Type-Options
- **Rate Limiting**: Token bucket algorithm with configurable RPS
- **Request ID Tracking**: UUID-based for audit trails
- **JWT Authentication**: With expiry configuration and claims extraction

### Observability
- **Structured Logging**: Zap logger with request ID correlation
- **Health Checks**: `/health` (liveness) and `/ready` (readiness with DB check)
- **Structured Errors**: Error codes and messages for API consumers

### Architecture
- **Clean Architecture**: Domain → Service → Handler layers
- **PostgreSQL + GORM**: With soft deletes and audit fields
- **Graceful Shutdown**: Signal handling for zero-downtime deployments
- **Input Validation**: go-playground/validator with structured error messages

## Project Structure

```text
├── cmd/api/              # Application entry point
├── internal/
│   ├── app/              # Server setup and DI
│   ├── config/           # Configuration management
│   ├── domain/           # Entities and repository interfaces
│   ├── dto/              # Request/Response data transfer objects
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # Security, rate limiting, auth, recovery
│   ├── repository/       # Data access implementations
│   └── service/          # Business logic
├── pkg/
│   ├── apperror/         # Structured error handling
│   ├── logger/           # Logging with context
│   ├── validator/        # Input validation
│   └── utils/fileutil/   # CSV, XLSX, PDF generation
├── migrations/           # Database migrations
└── .github/workflows/    # CI/CD pipeline
```

## Getting Started

### Prerequisites
- Go 1.22+
- PostgreSQL 14+

### Installation
```bash
# Clone and enter directory
cd go-boilerplate

# Copy environment file
cp .env.example .env

# Install dependencies
go mod tidy

# Start PostgreSQL (Docker example)
docker run -d --name postgres -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=boilerplate \
  postgres:14

# Run the application
go run cmd/api/main.go
```

## API Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/health` | Liveness probe | No |
| GET | `/ready` | Readiness probe | No |
| POST | `/auth/login` | Login and get JWT | No |
| POST | `/auth/register` | Register new user | No |
| GET | `/api/export/csv` | Export data as CSV | Yes |
| GET | `/api/export/xlsx` | Export data as XLSX | Yes |
| GET | `/api/export/pdf` | Export data as PDF | Yes |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_PORT` | 8080 | Server port |
| `APP_ENV` | development | Environment (development/production) |
| `JWT_SECRET` | - | JWT signing secret |
| `JWT_EXPIRY_HOURS` | 72 | Token expiry in hours |
| `RATE_LIMIT_RPS` | 10 | Requests per second limit |
| `RATE_LIMIT_BURST` | 20 | Burst size for rate limiter |
| `DB_*` | - | PostgreSQL connection details |

## License
MIT
