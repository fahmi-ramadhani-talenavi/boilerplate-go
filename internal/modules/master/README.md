# Master Module

Master/reference data management.

## Structure
```
master/
├── entity/         # Master data entities
├── handler/        # HTTP handlers per entity
├── migrations/     # 90+ mst_* tables
├── seeder/         # Master seeder logic
├── seeders/        # 29+ SQL seed files
└── module.go       # Module & routes setup
```

## Tables

| Table | Description |
|-------|-------------|
| mst_banks | Bank references |
| mst_branches | Branch offices |
| mst_provinces | Provinces |
| mst_districts | Districts/cities |
| mst_genders | Gender options |
| mst_religions | Religion options |
| mst_currencies | Currency types |
| mst_tax_brackets | Tax rates |
| ... | 40+ more tables |

## Endpoints

All require authentication.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/master/all?types=...` | **Batch Request** (multi-type) |
| GET | `/api/master/areas` | List areas |
| GET | `/api/master/banks` | List banks |
| GET | `/api/master/branches` | List branches |
| GET | `/api/master/provinces` | List provinces |
| GET | `/api/master/districts` | List districts |
| GET | `/api/master/genders` | List genders |
| GET | `/api/master/religions` | List religions |
| GET | `/api/master/currencies` | List currencies |
| GET | `/api/master/tax-groups`| List tax groups |
| GET | `/api/master/tax-brackets`| List tax brackets |

## Caching

The Batch API (`/all`) uses **Redis caching** to reduce database load.
- **Cache Key**: `master:<type>`
- **TTL**: 1 Hour
- **Behavior**: Automatically refreshes on cache miss.

## Seeding

Master data is seeded from SQL files in `seeders/` folder:
```bash
./cli seed:master
```
