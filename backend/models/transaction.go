package models

import (
	"time"
)

// Transaction represents a faucet transaction
type Transaction struct {
	ID            int64     `json:"id"`
	WalletAddress string    `json:"walletAddress"`
	IPAddress     string    `json:"ipAddress,omitempty"` // omitted in JSON responses
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"` // "pending", "completed", "failed"
	TxHash        string    `json:"txHash,omitempty"`
	ErrorMessage  string    `json:"errorMessage,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}

// ClaimHistory represents a user's claim history
type ClaimHistory struct {
	ID            int64     `json:"id"`
	WalletAddress string    `json:"walletAddress"`
	IPAddress     string    `json:"ipAddress,omitempty"` // omitted in JSON responses
	LastClaimTime time.Time `json:"lastClaimTime"`
	ClaimCount    int       `json:"claimCount"`
}

// CanClaim checks if a wallet can claim tokens based on the cooldown period
func (ch *ClaimHistory) CanClaim(cooldownSeconds int) (bool, time.Time) {
	nextClaimTime := ch.LastClaimTime.Add(time.Duration(cooldownSeconds) * time.Second)
	return time.Now().After(nextClaimTime), nextClaimTime
}

// TransactionResponse represents a transaction response to the client
type TransactionResponse struct {
	Success       bool         `json:"success"`
	Transaction   *Transaction `json:"transaction,omitempty"`
	Message       string       `json:"message,omitempty"`
	NextClaimTime time.Time    `json:"nextClaimTime,omitempty"`
}
