package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/maestroi/solana-faucet/backend/config"
	"github.com/maestroi/solana-faucet/backend/db"
	"github.com/maestroi/solana-faucet/backend/utils"
)

// Server represents the API server
type Server struct {
	config    *config.Config
	db        *db.Database
	router    *chi.Mux
	solana    *utils.SolanaClient
	server    *http.Server
	turnstile *utils.TurnstileClient

	// Balance caching
	balanceMutex    sync.RWMutex
	cachedBalance   float64
	lastBalanceTime time.Time
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, database *db.Database) *Server {
	r := chi.NewRouter()

	// Set up middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Set up CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not readily apparent
	}))

	// Create Solana client
	solanaClient, err := utils.NewSolanaClient(cfg.Solana.RpcURL, cfg.Solana.FaucetWalletPath)
	if err != nil {
		log.Fatalf("Failed to create Solana client: %v", err)
	}

	// Create Turnstile client
	turnstileClient := utils.NewTurnstileClient(cfg.Security.TurnstileSecretKey)

	// Create server
	s := &Server{
		config:    cfg,
		db:        database,
		router:    r,
		solana:    solanaClient,
		turnstile: turnstileClient,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port),
			Handler: r,
		},
	}

	// Set up routes
	s.setupRoutes()

	return s
}

// setupRoutes sets up the API routes
func (s *Server) setupRoutes() {
	s.router.Group(func(r chi.Router) {
		// Apply CORS middleware
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   s.config.CORS.AllowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))

		// Health check
		r.Get("/api/health", s.handleHealth)

		// Request funds
		r.Post("/api/request-funds", s.handleRequestFunds)

		// Get transactions
		r.Get("/api/transactions", s.handleGetTransactions)

		// Get balance
		r.Get("/api/balance", s.handleGetBalance)
	})
}

// Start starts the API server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the API server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// handleHealth handles the health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
