package usecases

import (
	"context"
	"errors"

	"github.com/sahdoio/crawlly-core/internal/membership/domain"
	"github.com/sahdoio/crawlly-core/pkg/auth"
)

// AuthenticateUserInput represents the data needed for login
type AuthenticateUserInput struct {
	Email    string
	Password string
}

// AuthenticateUserOutput represents the result of authentication
type AuthenticateUserOutput struct {
	UserID   string
	Email    string
	Name     string
	APIKey   string
	IsActive bool
}

// AuthenticateUserUseCase handles user authentication business logic
type AuthenticateUserUseCase struct {
	userRepo domain.UserRepository
}

// NewAuthenticateUserUseCase creates a new authentication use case
func NewAuthenticateUserUseCase(userRepo domain.UserRepository) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{userRepo: userRepo}
}

// Execute performs user authentication
func (uc *AuthenticateUserUseCase) Execute(ctx context.Context, input AuthenticateUserInput) (*AuthenticateUserOutput, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if !auth.CheckPassword(input.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// Return authenticated user data
	return &AuthenticateUserOutput{
		UserID:   user.ID.String(),
		Email:    user.Email,
		Name:     user.Name,
		APIKey:   user.APIKey,
		IsActive: user.IsActive,
	}, nil
}
