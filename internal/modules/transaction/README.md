# Transaction Module

Transaction and batch processing.

## Structure
```
transaction/
├── migrations/     # 6 app_* tables
├── seeders/        # SQL seed files
├── seeder/         # Seeder logic
└── (handlers TBD)
```

## Tables

| Table | Description |
|-------|-------------|
| app_batchings | Batch processing records |
| app_giro_reconciliations | Giro reconciliation |
| app_giro_reconciliation_details | Reconciliation details |

## Status

⚠️ **In Development** - Handlers not yet implemented.

Migration and seeder are ready:
```bash
./cli migrate:transaction
./cli seed:transaction
```
