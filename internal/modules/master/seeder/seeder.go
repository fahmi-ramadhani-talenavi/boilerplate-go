package seeder

import (
	"os"
	"path/filepath"

	"github.com/user/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Seeder handles master module seeding.
type Seeder struct {
	db *gorm.DB
}

// New creates a new master seeder.
func New(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// Seed runs all master data seeders from SQL files.
func (s *Seeder) Seed() error {
	seedPath := "internal/modules/master/seeders"
	
	files, err := os.ReadDir(seedPath)
	if err != nil {
		logger.Log.Warn("No master seeders found", zap.Error(err))
		return nil
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		tableName := file.Name()[:len(file.Name())-4]
		
		var count int64
		s.db.Table(tableName).Count(&count)
		if count > 0 {
			continue
		}

		content, err := os.ReadFile(filepath.Join(seedPath, file.Name()))
		if err != nil {
			logger.Log.Error("Failed to read seeder", zap.String("file", file.Name()), zap.Error(err))
			continue
		}

		if err := s.db.Exec(string(content)).Error; err != nil {
			logger.Log.Error("Failed to execute seeder", zap.String("file", file.Name()), zap.Error(err))
			continue
		}

		logger.Log.Debug("Seeded", zap.String("table", tableName))
	}

	logger.Log.Info("Master data seeded")
	return nil
}
