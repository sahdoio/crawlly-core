package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/sahdoio/crawlly-core/internal/membership/domain"
	"github.com/sahdoio/crawlly-core/pkg/auth"
)

// RegisterUserInput represents the data needed to register a user
type RegisterUserInput struct {
	Email    string
	Name     string
	Password string
}

// RegisterUserOutput represents the result of registration
type RegisterUserOutput struct {
	UserID string
	Email  string
	Name   string
	APIKey string
}

// RegisterUserUseCase handles user registration business logic
type RegisterUserUseCase struct {
	userRepo domain.UserRepository
}

// NewRegisterUserUseCase creates a new registration use case
func NewRegisterUserUseCase(userRepo domain.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{userRepo: userRepo}
}

// Execute performs the user registration
func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	// Validate input
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, _ := uc.userRepo.FindByEmail(ctx, input.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash the password
	passwordHash, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create new user
	user := domain.NewUser(input.Email, input.Name, passwordHash)

	// Save to database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Return output
	return &RegisterUserOutput{
		UserID: user.ID.String(),
		Email:  user.Email,
		Name:   user.Name,
		APIKey: user.APIKey,
	}, nil
}

// validateInput checks if the registration input is valid
func (uc *RegisterUserUseCase) validateInput(input RegisterUserInput) error {
	if strings.TrimSpace(input.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(input.Name) == "" {
		return errors.New("name is required")
	}
	if len(input.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	// Basic email validation
	if !strings.Contains(input.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}
