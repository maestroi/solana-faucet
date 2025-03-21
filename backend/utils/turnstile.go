package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// TurnstileClient is a client for interacting with the Cloudflare Turnstile API
type TurnstileClient struct {
	secretKey  string
	httpClient *http.Client
}

// TurnstileResponse represents a response from the Turnstile API
type TurnstileResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

// NewTurnstileClient creates a new Turnstile client
func NewTurnstileClient(secretKey string) *TurnstileClient {
	return &TurnstileClient{
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// VerifyToken verifies a Turnstile token
func (t *TurnstileClient) VerifyToken(token string) (bool, error) {
	// If no secret key is set, bypass verification (for development/testing)
	if t.secretKey == "" || t.secretKey == "your-turnstile-secret-key" {
		return true, nil
	}

	// Make request to Turnstile API
	resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", url.Values{
		"secret":   {t.secretKey},
		"response": {token},
	})
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Parse response
	var turnstileResp TurnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&turnstileResp); err != nil {
		return false, err
	}

	if !turnstileResp.Success && len(turnstileResp.ErrorCodes) > 0 {
		return false, fmt.Errorf("turnstile verification failed: %v", turnstileResp.ErrorCodes)
	}

	return turnstileResp.Success, nil
}
