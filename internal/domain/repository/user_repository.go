// Package repository defines the data access layer interfaces.
// Repository interfaces provide an abstraction over data storage,
// enabling the application to work with different databases without changes to business logic.
//
// DESIGN PRINCIPLE: The domain layer defines interfaces; the infrastructure layer implements them.
// This follows the Dependency Inversion Principle (DIP) from SOLID.
package repository

import (
	"context"

	"github.com/user/go-boilerplate/internal/domain/entity"
)

// UserRepository defines the contract for user data access operations.
// Implementations of this interface handle the actual database interactions.
//
// IMPLEMENTATIONS:
// - PostgresUserRepository: Production implementation using GORM
// - MockUserRepository: Test implementation for unit testing
//
// DESIGN NOTES:
// - All methods accept context for cancellation and tracing
// - Methods return domain entities, not DTOs
// - Errors should be wrapped with appropriate context
//
// SECURITY:
// - Implementations must use parameterized queries to prevent SQL injection
// - Sensitive data (passwords) should be handled carefully
type UserRepository interface {
	// GetByID retrieves a user by their unique identifier.
	//
	// PARAMETERS:
	// - ctx: Context for cancellation and request tracing
	// - id: UUID of the user to retrieve
	//
	// RETURNS:
	// - *entity.User: User entity if found
	// - error: gorm.ErrRecordNotFound if not found, or other database error
	//
	// USE CASE: Loading user profile, authorization checks
	GetByID(ctx context.Context, id string) (*entity.User, error)

	// GetByEmail retrieves a user by their email address.
	//
	// PARAMETERS:
	// - ctx: Context for cancellation and request tracing
	// - email: Email address to search for (case-sensitive)
	//
	// RETURNS:
	// - *entity.User: User entity if found
	// - error: gorm.ErrRecordNotFound if not found, or other database error
	//
	// USE CASE: Authentication (login), duplicate email checking
	//
	// SECURITY: Email lookups should be indexed for performance
	// to prevent timing attacks during authentication.
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Create inserts a new user record into the database.
	//
	// PARAMETERS:
	// - ctx: Context for cancellation and request tracing
	// - user: User entity to create (ID will be auto-generated if empty)
	//
	// RETURNS:
	// - error: Database error if creation fails (e.g., duplicate email)
	//
	// USE CASE: User registration
	//
	// SECURITY:
	// - Password must be hashed before calling this method
	// - Email uniqueness is enforced by database constraint
	Create(ctx context.Context, user *entity.User) error

	// Update modifies an existing user record.
	//
	// PARAMETERS:
	// - ctx: Context for cancellation and request tracing
	// - user: User entity with updated fields (ID must be set)
	//
	// RETURNS:
	// - error: Database error if update fails
	//
	// USE CASE: Profile updates, password changes, account status changes
	//
	// AUDIT: Updates are tracked via UpdatedAt and UpdatedBy fields
	Update(ctx context.Context, user *entity.User) error

	// Delete soft-deletes a user record.
	// The record is not permanently removed; DeletedAt is set instead.
	//
	// PARAMETERS:
	// - ctx: Context for cancellation and request tracing
	// - id: UUID of the user to delete
	//
	// RETURNS:
	// - error: Database error if deletion fails
	//
	// USE CASE: Account deletion requests
	//
	// COMPLIANCE: Soft delete maintains data for audit and recovery purposes.
	// Consider data retention policies for GDPR compliance.
	Delete(ctx context.Context, id string) error

	// List retrieves a paginated list of users.
	//
	// PARAMETERS:
	// - ctx: Context for cancellation and request tracing
	// - offset: Number of records to skip (for pagination)
	// - limit: Maximum number of records to return
	//
	// RETURNS:
	// - []*entity.User: List of user entities
	// - int64: Total count of users (for pagination metadata)
	// - error: Database error if query fails
	//
	// USE CASE: Admin user management, user listing
	//
	// SECURITY: This endpoint should be protected and only accessible
	// by administrators.
	List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error)
}
