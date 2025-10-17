package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
// This is an INTERFACE - it defines what operations are possible,
// but not HOW they're implemented. This allows us to swap
// implementations (PostgreSQL, MySQL, MongoDB, etc.) easily.
type UserRepository interface {
	// Create saves a new user to the database
	Create(ctx context.Context, user *User) error

	// FindByID retrieves a user by their ID
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	// FindByEmail retrieves a user by their email address
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindByAPIKey retrieves a user by their API key
	FindByAPIKey(ctx context.Context, apiKey string) (*User, error)

	// Update saves changes to an existing user
	Update(ctx context.Context, user *User) error

	// Delete removes a user from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all users with pagination
	// offset: how many to skip, limit: how many to return
	List(ctx context.Context, offset, limit int) ([]*User, error)
}