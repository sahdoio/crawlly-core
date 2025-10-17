package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/sahdoio/crawlly-core/internal/membership/handlers"
	"github.com/sahdoio/crawlly-core/internal/membership/repositories"
	"github.com/sahdoio/crawlly-core/internal/membership/usecases"
	"github.com/sahdoio/crawlly-core/pkg/config"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connected successfully")

	// Initialize repositories
	userRepo := repositories.NewPostgresUserRepository(db)

	// Initialize use cases
	registerUseCase := usecases.NewRegisterUserUseCase(userRepo)
	authenticateUseCase := usecases.NewAuthenticateUserUseCase(userRepo)

	// Initialize handlers
	authHandlers := handlers.NewAuthHandlers(registerUseCase, authenticateUseCase)

	// Initialize router (chi is a lightweight HTTP router)
	r := chi.NewRouter()

	// Add middleware (functions that run on every request)
	r.Use(middleware.Logger)      // Logs all requests
	r.Use(middleware.Recoverer)   // Recovers from panics
	r.Use(middleware.RequestID)   // Adds unique ID to each request

	// Define a simple health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Authentication routes
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandlers.Register)
		r.Post("/login", authHandlers.Login)
	})

	// Create HTTP/2 server with h2c (HTTP/2 Cleartext - without TLS)
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      h2c.NewHandler(r, &http2.Server{}),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine (like a thread)
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Blocks until signal received

	log.Println("Shutting down server...")

	// Give server 30 seconds to finish existing requests
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}