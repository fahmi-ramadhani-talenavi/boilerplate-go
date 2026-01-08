// Package entity contains the domain entities for the application.
package entity

// User represents a user account in the system.
// This entity is central to authentication and authorization.
//
// SECURITY CONSIDERATIONS:
// - Password is excluded from JSON serialization to prevent exposure
// - IsActive flag enables account suspension without deletion
// - Email has unique constraint to prevent duplicate accounts
//
// COMPLIANCE:
// - Inherits audit fields from Base for regulatory requirements
// - Soft delete ensures data retention policies are met
type User struct {
	// Base embeds common fields: ID, CreatedAt, UpdatedAt, DeletedAt, CreatedBy, UpdatedBy
	Base

	// Email is the user's unique email address.
	// Used as the primary identifier for authentication.
	// Indexed for fast lookups during login.
	//
	// VALIDATION: Must be a valid email format (enforced in DTO layer)
	Email string `json:"email" gorm:"uniqueIndex;not null"`

	// Password stores the bcrypt-hashed password.
	// NEVER stored in plain text.
	//
	// SECURITY:
	// - Excluded from JSON to prevent accidental exposure
	// - Uses bcrypt with cost factor for brute-force resistance
	// - Minimum 8 characters enforced at validation layer
	Password string `json:"-" gorm:"not null"`

	// Name is the user's display name.
	// Used for personalization in the UI.
	Name string `json:"name"`

	// IsActive indicates whether the account is enabled.
	// Set to false to suspend an account without deleting it.
	//
	// USE CASES:
	// - Suspicious activity detected
	// - Failed payment/subscription
	// - User-requested temporary deactivation
	// - Administrative action
	IsActive bool `json:"is_active" gorm:"default:true"`
}

// TableName specifies the database table name for GORM.
// This overrides the default pluralized name.
func (User) TableName() string {
	return "users"
}
