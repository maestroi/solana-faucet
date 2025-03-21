package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/maestroi/solana-faucet/backend/models"
	"github.com/maestroi/solana-faucet/backend/utils"
)

// BalanceResponse represents the response for the balance endpoint
type BalanceResponse struct {
	Balance float64 `json:"balance"`
	Cached  bool    `json:"cached"`
}

// Cache duration for balance
const balanceCacheDuration = 1 * time.Minute

// handleGetBalance returns the current balance of the faucet wallet
func (s *Server) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Balance] Starting balance request from %s", r.RemoteAddr)

	s.balanceMutex.RLock()
	// Check if we have a cached balance that's less than 1 minute old
	if !s.lastBalanceTime.IsZero() && time.Since(s.lastBalanceTime) < balanceCacheDuration {
		balance := s.cachedBalance
		s.balanceMutex.RUnlock()
		log.Printf("[Balance] Returning cached balance: %f SOL", balance)

		response := BalanceResponse{
			Balance: balance,
			Cached:  true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	s.balanceMutex.RUnlock()

	// Get fresh balance
	s.balanceMutex.Lock()
	defer s.balanceMutex.Unlock()

	// Double check if another request already updated the cache
	if !s.lastBalanceTime.IsZero() && time.Since(s.lastBalanceTime) < balanceCacheDuration {
		log.Printf("[Balance] Another request updated cache, returning cached balance: %f SOL", s.cachedBalance)
		response := BalanceResponse{
			Balance: s.cachedBalance,
			Cached:  true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("[Balance] Fetching fresh balance from Solana")
	balance, err := s.solana.GetFaucetBalance()
	if err != nil {
		log.Printf("[Balance] Error getting balance: %v", err)
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}
	log.Printf("[Balance] Got raw balance in lamports: %d", balance)

	// Convert lamports to SOL (1 SOL = 1e9 lamports)
	balanceSOL := float64(balance) / 1e9

	// Update cache
	s.cachedBalance = balanceSOL
	s.lastBalanceTime = time.Now()

	response := BalanceResponse{
		Balance: balanceSOL,
		Cached:  false,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("[Balance] Successfully returned fresh balance")
}

// handleRequestFunds handles the request funds endpoint
func (s *Server) handleRequestFunds(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req models.FundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[RequestFunds] Invalid request body: %v", err)
		response := map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Log the request for debugging
	log.Printf("[RequestFunds] Received request for wallet: %s", req.WalletAddress)

	// Validate required fields
	if req.WalletAddress == "" {
		response := map[string]interface{}{
			"success": false,
			"error":   "Wallet address is required",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.TurnstileResponse == "" {
		response := map[string]interface{}{
			"success": false,
			"error":   "Turnstile response is required",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get client IP
	clientIP := r.RemoteAddr
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		clientIP = forwardedFor
	}

	// Validate Turnstile token
	if s.turnstile != nil {
		isValid, err := s.turnstile.VerifyToken(req.TurnstileResponse)
		if err != nil {
			log.Printf("[RequestFunds] Turnstile verification error: %v", err)
			response := map[string]interface{}{
				"success": false,
				"error":   "Failed to verify Turnstile token",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		if !isValid {
			response := map[string]interface{}{
				"success": false,
				"error":   "Invalid Turnstile token",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Validate wallet address
	if !utils.IsValidSolanaAddress(req.WalletAddress) {
		response := map[string]interface{}{
			"success": false,
			"error":   "Invalid Solana wallet address format",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if this wallet has claimed recently
	history, err := s.db.GetClaimHistory(req.WalletAddress)
	if err != nil {
		log.Printf("[RequestFunds] Error checking claim history: %v", err)
		response := map[string]interface{}{
			"success": false,
			"error":   "Failed to check claim history",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if history != nil {
		// Check if the wallet can claim again
		canClaim, nextClaimTime := history.CanClaim(s.config.Security.ClaimCooldown)
		if !canClaim {
			// Calculate wait time
			waitTime := time.Until(nextClaimTime)
			hours := int(waitTime.Hours())
			minutes := int(waitTime.Minutes()) % 60
			seconds := int(waitTime.Seconds()) % 60

			var timeMsg string
			if hours >= 2 {
				// Format the next claim time as a date
				timeMsg = nextClaimTime.Format("Jan 2 at 3:04 PM")
			} else if hours > 0 {
				if minutes > 0 {
					timeMsg = fmt.Sprintf("%d hour%s and %d minute%s",
						hours, pluralize(hours),
						minutes, pluralize(minutes))
				} else {
					timeMsg = fmt.Sprintf("%d hour%s", hours, pluralize(hours))
				}
			} else if minutes > 0 {
				if seconds > 0 {
					timeMsg = fmt.Sprintf("%d minute%s and %d second%s",
						minutes, pluralize(minutes),
						seconds, pluralize(seconds))
				} else {
					timeMsg = fmt.Sprintf("%d minute%s", minutes, pluralize(minutes))
				}
			} else {
				timeMsg = fmt.Sprintf("%d second%s", seconds, pluralize(seconds))
			}

			response := map[string]interface{}{
				"success":       false,
				"error":         fmt.Sprintf("Please wait until %s before requesting funds again", timeMsg),
				"nextClaimTime": nextClaimTime.Unix(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Send transaction
	txHash, err := s.solana.SendSOL(req.WalletAddress, s.config.Solana.AmountPerRequest)
	if err != nil {
		log.Printf("[RequestFunds] Error sending transaction: %v", err)
		response := map[string]interface{}{
			"success": false,
			"error":   "Failed to send transaction",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Record transaction in database
	tx := &models.Transaction{
		WalletAddress: req.WalletAddress,
		IPAddress:     clientIP,
		Amount:        s.config.Solana.AmountPerRequest,
		Status:        "completed",
		TxHash:        txHash,
		Timestamp:     time.Now(),
	}
	if _, err := s.db.CreateTransaction(tx); err != nil {
		log.Printf("[RequestFunds] Failed to save transaction: %v", err)
	}

	// Update claim history
	if err := s.db.UpdateClaimHistory(req.WalletAddress, clientIP); err != nil {
		log.Printf("[RequestFunds] Failed to update claim history: %v", err)
	}

	// Return success response
	response := map[string]interface{}{
		"success":          true,
		"amount":           s.config.Solana.AmountPerRequest,
		"transaction_hash": txHash,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to add 's' for plurals
func pluralize(n int) string {
	if n != 1 {
		return "s"
	}
	return ""
}

// handleGetTransactions returns the recent transactions
func (s *Server) handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := s.db.GetRecentTransactions(10)
	if err != nil {
		log.Printf("[Transactions] Error getting transactions: %v", err)
		response := models.TransactionResponse{
			Success: false,
			Message: "Failed to get transactions",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := struct {
		Success      bool                  `json:"success"`
		Transactions []*models.Transaction `json:"transactions"`
	}{
		Success:      true,
		Transactions: transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
