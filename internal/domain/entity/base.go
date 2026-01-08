// Package entity contains the domain entities for the application.
// These entities represent the core business objects and are independent
// of any infrastructure concerns like databases or HTTP.
package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base provides common fields for all database entities.
// It implements audit trail functionality required for banking-grade applications:
// - ID: Unique identifier using UUID for security (no sequential IDs)
// - CreatedAt/UpdatedAt: Automatic timestamps for audit trails
// - DeletedAt: Soft delete support (data is never permanently deleted)
// - CreatedBy/UpdatedBy: User attribution for compliance and auditing
//
// SECURITY NOTE: Never expose DeletedAt in API responses.
// COMPLIANCE: All modifications are tracked for regulatory requirements.
type Base struct {
	// ID is the primary key using UUID v4 for unpredictability.
	// Using UUID instead of auto-increment prevents enumeration attacks.
	ID string `json:"id" gorm:"primaryKey;type:uuid"`

	// CreatedAt records when the record was first created.
	// This field is automatically managed by GORM.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt records the last modification time.
	// This field is automatically managed by GORM.
	UpdatedAt time.Time `json:"updated_at"`

	// DeletedAt enables soft delete functionality.
	// Records are never permanently removed for audit compliance.
	// Hidden from JSON responses for security.
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// CreatedBy stores the ID of the user who created this record.
	// Required for audit trails and compliance reporting.
	// Using pointer to allow NULL for system-generated records.
	CreatedBy *string `json:"created_by,omitempty" gorm:"type:uuid"`

	// UpdatedBy stores the ID of the user who last modified this record.
	// Required for audit trails and compliance reporting.
	// Using pointer to allow NULL for system-generated records.
	UpdatedBy *string `json:"updated_by,omitempty" gorm:"type:uuid"`
}

// BeforeCreate is a GORM hook that runs before inserting a new record.
// It automatically generates a UUID if one is not provided.
//
// SECURITY: Uses crypto-secure UUID v4 for unpredictable identifiers.
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
