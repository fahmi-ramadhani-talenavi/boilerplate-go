# System Module

System configuration and settings.

## Structure
```
system/
├── entity/         # System entities
├── handler/        # API handlers
├── migrations/     # 24+ sys_* tables
├── seeder/         # Seeder logic
├── seeders/        # SQL seed files
└── module.go       # Module & routes setup
```

## Tables

| Table | Description |
|-------|-------------|
| sys_settings | Application settings |
| sys_roles | User roles |
| sys_sub_roles | Sub-roles |
| sys_bank_fees | Bank fee config |
| sys_base_fees | Base fee config |
| sys_sub_menus | Menu structure |
| sys_announcements | Announcements |
| ... | 15+ more tables |

## Endpoints

All require authentication.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/system/settings` | Get all settings (key-value) |
| GET | `/api/system/roles` | List roles |
| GET | `/api/system/sub-roles` | List sub-roles |
| GET | `/api/system/bank-fees` | List bank fees |
| GET | `/api/system/base-fees` | List base fees |
| GET | `/api/system/menus` | List menu structure |

## Seeding

```bash
./cli seed:system
```
