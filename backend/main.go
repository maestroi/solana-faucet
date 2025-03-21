package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/maestroi/solana-faucet/backend/api"
	"github.com/maestroi/solana-faucet/backend/config"
	"github.com/maestroi/solana-faucet/backend/db"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize database
	database, err := db.InitDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer database.Close()

	// Set up API server
	server := api.NewServer(cfg, database)

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on %s", cfg.Server.Address)
		if err := server.Start(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")
	if err := server.Shutdown(); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	fmt.Println("Server stopped")
}
