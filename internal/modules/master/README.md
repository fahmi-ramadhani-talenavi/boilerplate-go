# Master Module

Master/reference data management.

## Structure
```
master/
├── migrations/     # 90 mst_* tables
├── seeders/        # 29 SQL seed files
├── seeder/         # Seeder logic
├── entity.go       # Master data entities
├── handler.go      # List handlers
└── module.go       # Routes setup
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
| GET | `/api/master/banks` | List banks |
| GET | `/api/master/branches` | List branches |
| GET | `/api/master/provinces` | List provinces |
| GET | `/api/master/districts` | List districts |
| GET | `/api/master/genders` | List genders |
| GET | `/api/master/religions` | List religions |
| GET | `/api/master/currencies` | List currencies |

## Seeding

Master data is seeded from SQL files in `seeders/` folder:
```bash
./cli seed:master
```
