package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sahdoio/crawlly-core/internal/membership/usecases"
)

// AuthHandlers contains all authentication-related HTTP handlers
type AuthHandlers struct {
	registerUseCase     *usecases.RegisterUserUseCase
	authenticateUseCase *usecases.AuthenticateUserUseCase
}

// NewAuthHandlers creates a new authentication handlers instance
func NewAuthHandlers(
	registerUseCase *usecases.RegisterUserUseCase,
	authenticateUseCase *usecases.AuthenticateUserUseCase,
) *AuthHandlers {
	return &AuthHandlers{
		registerUseCase:     registerUseCase,
		authenticateUseCase: authenticateUseCase,
	}
}

// RegisterRequest represents the JSON body for registration
type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// RegisterResponse represents the JSON response for registration
type RegisterResponse struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	APIKey  string `json:"api_key"`
	Message string `json:"message"`
}

// LoginRequest represents the JSON body for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the JSON response for login
type LoginResponse struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	APIKey  string `json:"api_key"`
	Message string `json:"message"`
}

// ErrorResponse represents a JSON error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Register handles POST /api/auth/register
func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Execute use case
	input := usecases.RegisterUserInput{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	output, err := h.registerUseCase.Execute(r.Context(), input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Respond with success
	respondWithJSON(w, http.StatusCreated, RegisterResponse{
		UserID:  output.UserID,
		Email:   output.Email,
		Name:    output.Name,
		APIKey:  output.APIKey,
		Message: "User registered successfully",
	})
}

// Login handles POST /api/auth/login
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Execute use case
	input := usecases.AuthenticateUserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.authenticateUseCase.Execute(r.Context(), input)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Respond with success
	respondWithJSON(w, http.StatusOK, LoginResponse{
		UserID:  output.UserID,
		Email:   output.Email,
		Name:    output.Name,
		APIKey:  output.APIKey,
		Message: "Login successful",
	})
}

// respondWithJSON writes a JSON response
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// respondWithError writes a JSON error response
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, ErrorResponse{Error: message})
}
