package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered user in the system
type User struct {
	ID           uuid.UUID  // Unique identifier (UUID is better than int for distributed systems)
	Email        string     // User's email (must be unique)
	Name         string     // Display name
	PasswordHash string     // Bcrypt hashed password (never store plain passwords!)
	APIKey       string     // API key for authenticating crawler requests
	IsActive     bool       // Whether the account is active
	CreatedAt    time.Time  // When the account was created
	UpdatedAt    time.Time  // Last time the account was updated
}

// NewUser creates a new User with generated UUID and timestamps
func NewUser(email, name, passwordHash string) *User {
	now := time.Now()
	return &User{
		ID:           uuid.New(),
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
		APIKey:       generateAPIKey(),
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// generateAPIKey creates a unique API key for the user
func generateAPIKey() string {
	// Generate a UUID and use it as API key
	// In production, you might want a more sophisticated approach
	return "crawlly_" + uuid.New().String()
}

// UpdateEmail updates the user's email and timestamp
func (u *User) UpdateEmail(email string) {
	u.Email = email
	u.UpdatedAt = time.Now()
}

// UpdateName updates the user's name and timestamp
func (u *User) UpdateName(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}

// Deactivate marks the user account as inactive
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate marks the user account as active
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// RegenerateAPIKey creates a new API key for the user
func (u *User) RegenerateAPIKey() {
	u.APIKey = generateAPIKey()
	u.UpdatedAt = time.Now()
}