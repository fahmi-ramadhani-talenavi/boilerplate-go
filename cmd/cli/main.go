package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/user/go-boilerplate/internal/config"
	"github.com/user/go-boilerplate/internal/database/migration"
	authseeder "github.com/user/go-boilerplate/internal/modules/auth/seeder"
	masterseeder "github.com/user/go-boilerplate/internal/modules/master/seeder"
	systemseeder "github.com/user/go-boilerplate/internal/modules/system/seeder"
	transactionseeder "github.com/user/go-boilerplate/internal/modules/transaction/seeder"
	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var moduleMigrations = map[string]string{
	"auth":        "file://internal/modules/auth/migrations",
	"master":      "file://internal/modules/master/migrations",
	"system":      "file://internal/modules/system/migrations",
	"transaction": "file://internal/modules/transaction/migrations",
}

var migrationOrder = []string{"auth", "system", "master", "transaction"}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	cfg := config.LoadConfig()
	logger.InitLogger(cfg.LogLevel, cfg.AppEnv)
	defer logger.Log.Sync()

	switch command {
	case "migrate":
		for _, mod := range migrationOrder {
			db_loop, err := initDatabase(cfg)
			if err != nil {
				logger.Log.Fatal("Failed to connect to database", zap.Error(err))
			}
			sqlDB_loop, _ := db_loop.DB()
			
			path := moduleMigrations[mod]
			tableName := "schema_migrations_" + mod
			fmt.Printf("-> Migrating: %s\n", mod)
			if err := runMigration(sqlDB_loop, path, tableName); err != nil {
				logger.Log.Fatal("Migration failed", zap.String("module", mod), zap.Error(err))
			}
			sqlDB_loop.Close()
		}
		fmt.Println("V All migrations completed")

	case "migrate:auth", "migrate:master", "migrate:system", "migrate:transaction":
		mod := command[8:]
		db_loop, _ := initDatabase(cfg)
		sqlDB_loop, _ := db_loop.DB()
		path := moduleMigrations[mod]
		tableName := "schema_migrations_" + mod
		fmt.Printf("-> Migrating: %s\n", mod)
		if err := runMigration(sqlDB_loop, path, tableName); err != nil {
			logger.Log.Fatal("Migration failed", zap.Error(err))
		}
		sqlDB_loop.Close()
		fmt.Printf("V %s migrated\n", mod)

	case "migrate:rollback":
		count := 1
		if len(os.Args) > 2 {
			if n, err := strconv.Atoi(os.Args[2]); err == nil {
				count = n
			}
		}
		for i := len(migrationOrder) - 1; i >= 0; i-- {
			mod := migrationOrder[i]
			db_loop, _ := initDatabase(cfg)
			sqlDB_loop, _ := db_loop.DB()
			path := moduleMigrations[mod]
			tableName := "schema_migrations_" + mod
			fmt.Printf("-> Rolling back: %s\n", mod)
			runRollback(sqlDB_loop, path, tableName, count)
			sqlDB_loop.Close()
		}
		fmt.Println("V Rollback completed")

	case "migrate:status":
		for _, mod := range migrationOrder {
			db_loop, _ := initDatabase(cfg)
			sqlDB_loop, _ := db_loop.DB()
			path := moduleMigrations[mod]
			tableName := "schema_migrations_" + mod
			fmt.Printf("\n=== %s ===\n", mod)
			showStatus(sqlDB_loop, path, tableName)
			sqlDB_loop.Close()
		}

	case "migrate:fresh":
		fmt.Println("! Dropping all tables...")
		for i := len(migrationOrder) - 1; i >= 0; i-- {
			db_loop, _ := initDatabase(cfg)
			sqlDB_loop, _ := db_loop.DB()
			path := moduleMigrations[migrationOrder[i]]
			tableName := "schema_migrations_" + migrationOrder[i]
			if m, err := migration.NewMigrator(sqlDB_loop, path, tableName); err == nil {
				m.Drop()
			}
			sqlDB_loop.Close()
		}
		fmt.Println("! Re-migrating all modules...")
		for _, mod := range migrationOrder {
			db_loop, _ := initDatabase(cfg)
			sqlDB_loop, _ := db_loop.DB()
			path := moduleMigrations[mod]
			tableName := "schema_migrations_" + mod
			fmt.Printf("-> Migrating: %s\n", mod)
			runMigration(sqlDB_loop, path, tableName)
			sqlDB_loop.Close()
		}
		fmt.Println("V Database reset completed")

	case "seed":
		db_seed, _ := initDatabase(cfg)
		fmt.Println("-> Seeding auth...")
		authseeder.New(db_seed).Seed()
		fmt.Println("-> Seeding system...")
		systemseeder.New(db_seed).Seed()
		fmt.Println("-> Seeding master...")
		masterseeder.New(db_seed).Seed()
		fmt.Println("-> Seeding transaction...")
		transactionseeder.New(db_seed).Seed()
		fmt.Println("V All seeders completed")

	case "seed:auth":
		db_seed, _ := initDatabase(cfg)
		authseeder.New(db_seed).Seed()
		fmt.Println("V Auth seeded")

	case "seed:master":
		db_seed, _ := initDatabase(cfg)
		masterseeder.New(db_seed).Seed()
		fmt.Println("V Master seeded")

	case "seed:system":
		db_seed, _ := initDatabase(cfg)
		systemseeder.New(db_seed).Seed()
		fmt.Println("V System seeded")

	case "seed:transaction":
		db_seed, _ := initDatabase(cfg)
		transactionseeder.New(db_seed).Seed()
		fmt.Println("V Transaction seeded")

	default:
		printUsage()
		os.Exit(1)
	}
}

func runMigration(sqlDB *sql.DB, path string, tableName string) error {
	m, err := migration.NewMigrator(sqlDB, path, tableName)
	if err != nil {
		return err
	}
	return m.Up()
}

func runRollback(sqlDB *sql.DB, path string, tableName string, count int) error {
	m, err := migration.NewMigrator(sqlDB, path, tableName)
	if err != nil {
		return err
	}
	return m.Down(count)
}

func showStatus(sqlDB *sql.DB, path string, tableName string) {
	m, err := migration.NewMigrator(sqlDB, path, tableName)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
		return
	}
	m.Status()
}

func printUsage() {
	fmt.Println(`
Database CLI

Migrations:
  migrate              All modules
  migrate:auth         Auth only
  migrate:master       Master only
  migrate:system       System only
  migrate:transaction  Transaction only
  migrate:rollback     Rollback
  migrate:status       Status
  migrate:fresh        Drop & re-migrate

Seeders:
  seed                 All modules
  seed:auth            Auth only
  seed:master          Master only
  seed:system          System only
  seed:transaction     Transaction only
`)
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	gormConfig := &gorm.Config{}
	if cfg.AppEnv == "production" {
		gormConfig.Logger = gormlogger.Default.LogMode(gormlogger.Silent)
	} else {
		gormConfig.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
