# Auth Module

Authentication and user management.

## Structure
```
auth/
├── dto/            # Data Transfer Objects
├── entity/         # Database entities
├── handler/        # HTTP handlers
├── migrations/     # sys_users table
├── repository/     # Data access layer
├── seeder/         # Seeder logic
├── seeders/        # SQL seed files
├── service/        # Business logic
└── module.go       # Module & routes setup
```

## Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/login` | Login, returns JWT |
| POST | `/auth/register` | Register new user |
| GET | `/auth/me` | Get current user (auth required) |

## Environment

| Variable | Description |
|----------|-------------|
| `JWT_SECRET` | JWT signing key |
| `JWT_EXPIRY_HOURS` | Token expiry (default: 72) |
| `ADMIN_PASSWORD` | Default admin password |

## Default Users (Seeder)

| Email | Password | Notes |
|-------|----------|-------|
| admin@example.com | Admin@123456 | Change in production! |
