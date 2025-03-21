package config

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		Address string
		Port    int
	}
	Database struct {
		Path string
	}
	Solana struct {
		RpcURL             string
		FaucetWalletPath   string
		AmountPerRequest   float64
		NetworkType        string // "testnet", "devnet", etc.
		TransactionTimeout int
	}
	Security struct {
		TurnstileSecretKey string
		TurnstileSiteKey   string
		RateLimitRequests  int
		RateLimitDuration  int // in seconds
		ClaimCooldown      int // in seconds
	}
	CORS struct {
		AllowedOrigins []string
	}
}

// LoadConfig loads the application configuration from environment variables
func LoadConfig(_ string) (*Config, error) {
	var config Config

	// Server config
	config.Server.Address = getEnvWithDefault("FAUCET_SERVER_ADDRESS", "0.0.0.0")
	config.Server.Port = getEnvIntWithDefault("FAUCET_SERVER_PORT", 8080)

	// Database config
	config.Database.Path = getEnvWithDefault("FAUCET_DB_PATH", "/app/data/faucet.db")

	// Solana config
	config.Solana.RpcURL = getEnvWithDefault("FAUCET_SOLANA_RPC_URL", "https://api.testnet.solana.com")
	config.Solana.FaucetWalletPath = getEnvWithDefault("FAUCET_WALLET_PATH", "/app/data/wallet.json")
	config.Solana.AmountPerRequest = getEnvFloatWithDefault("FAUCET_AMOUNT_PER_REQUEST", 1.0)
	config.Solana.NetworkType = getEnvWithDefault("FAUCET_NETWORK_TYPE", "testnet")
	config.Solana.TransactionTimeout = getEnvIntWithDefault("FAUCET_TRANSACTION_TIMEOUT", 30)

	// Security config
	config.Security.TurnstileSecretKey = getEnvWithDefault("FAUCET_TURNSTILE_SECRET", "your-turnstile-secret-key")
	config.Security.TurnstileSiteKey = getEnvWithDefault("FAUCET_TURNSTILE_SITE", "your-turnstile-site-key")
	config.Security.RateLimitRequests = getEnvIntWithDefault("FAUCET_RATE_LIMIT_REQUESTS", 5)
	config.Security.RateLimitDuration = getEnvIntWithDefault("FAUCET_RATE_LIMIT_DURATION", 60)
	config.Security.ClaimCooldown = getEnvIntWithDefault("FAUCET_CLAIM_COOLDOWN", 86400)

	// CORS config
	allowedOrigins := getEnvWithDefault("FAUCET_CORS_ALLOWED_ORIGINS", "http://localhost:3000,https://https://solana-faucet.maestroi.cc/")
	config.CORS.AllowedOrigins = strings.Split(allowedOrigins, ",")

	return &config, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloatWithDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// CreateDefaultConfig creates a default configuration file if one doesn't exist
func CreateDefaultConfig(path string) error {
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return nil // File exists, no need to create
	}

	// Create a default configuration
	config := Config{}
	config.Server.Address = "0.0.0.0"
	config.Server.Port = 8080
	config.Database.Path = "faucet.db"
	config.Solana.RpcURL = "https://api.testnet.solana.com"
	config.Solana.FaucetWalletPath = "wallet.json"
	config.Solana.AmountPerRequest = 1.0
	config.Solana.NetworkType = "testnet"
	config.Solana.TransactionTimeout = 30
	config.Security.TurnstileSecretKey = "your-turnstile-secret-key"
	config.Security.TurnstileSiteKey = "your-turnstile-site-key"
	config.Security.RateLimitRequests = 5
	config.Security.RateLimitDuration = 60
	config.Security.ClaimCooldown = 86400 // 24 hours in seconds

	// Create the file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the configuration to JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}
