// Package entity provides base entity structures for all modules.
package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base provides common fields for all database entities.
type Base struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	CreatedBy *string        `json:"created_by,omitempty" gorm:"type:uuid"`
	UpdatedBy *string        `json:"updated_by,omitempty" gorm:"type:uuid"`
}

// BeforeCreate generates UUID if not provided.
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
