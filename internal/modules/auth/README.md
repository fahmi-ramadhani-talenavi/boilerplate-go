# Auth Module

Authentication and user management.

## Structure
```
auth/
├── migrations/     # sys_users table
├── seeder/         # Admin & sample users
├── entity.go       # User entity
├── dto.go          # Login/Register DTOs
├── repository.go   # User repository
├── service.go      # Auth business logic
├── handler.go      # HTTP handlers
└── module.go       # Module setup
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
