# Go Boilerplate

Production-ready modular monolith for financial applications.

## Architecture

**Modular Monolith** with clean code structure:

```
cmd/
├── api/                    # HTTP Server binary
└── cli/                    # Migration & Seeder binary

internal/modules/
├── auth/
│   ├── entity/user.go
│   ├── dto/auth.go
│   ├── repository/user.go
│   ├── service/auth.go
│   ├── handler/auth.go
│   ├── seeder/
│   ├── migrations/
│   └── module.go
├── master/                 # 13 entities (bank, province, etc)
├── system/                 # 7 entities (role, settings, fees)
├── transaction/
├── file/
└── health/
```

## Build & Run

```bash
# Install dependencies
go mod tidy

# Build binaries
go build -o build/api ./cmd/api    # HTTP server only
go build -o build/cli ./cmd/cli    # Migration/seeder only

# Or build all at once (into root, not recommended for this setup)
# go build -o build/ ./cmd/...
```

## Database Setup

```bash
# Run all migrations (auth -> system -> master -> transaction)
./build/cli migrate

# Run specific module migration
./build/cli migrate:auth
./build/cli migrate:master
./build/cli migrate:system
./build/cli migrate:transaction

# Check migration status
./build/cli migrate:status

# Rollback last migration
./build/cli migrate:rollback

# Fresh database (drop all & re-migrate)
./build/cli migrate:fresh
```

## Seeding

```bash
# Run all seeders
./build/cli seed

# Run specific module seeder
./build/cli seed:auth           # Admin user
./build/cli seed:master         # Master data (banks, provinces, etc)
./build/cli seed:system         # System settings, roles, fees
./build/cli seed:transaction    # Transaction data
```

## Run Server

```bash
# Development
./build/api

# Or directly
go run ./cmd/api
```

Server starts at `http://localhost:8080`

## API Endpoints

| Endpoint | Auth | Description |
|----------|------|-------------|
| `GET /health` | ❌ | Liveness probe |
| `GET /ready` | ❌ | Readiness (DB check) |
| `POST /auth/login` | ❌ | Login |
| `POST /auth/register` | ❌ | Register |
| `GET /auth/me` | ✅ | Current user |
| `GET /api/master/*` | ✅ | Master data |
| `GET /api/system/*` | ✅ | System config |
| `POST /api/upload` | ✅ | File upload |

## Configuration

Create `.env` file:
```env
APP_PORT=8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=boilerplate
DB_SSLMODE=disable

JWT_SECRET=your-secret-key
JWT_EXPIRY_HOURS=72

REDIS_HOST=localhost
REDIS_PORT=6379
```

## Module Documentation

- [Auth Module](internal/modules/auth/README.md)
- [Master Module](internal/modules/master/README.md)
- [System Module](internal/modules/system/README.md)
- [Transaction Module](internal/modules/transaction/README.md)

## License
MIT
