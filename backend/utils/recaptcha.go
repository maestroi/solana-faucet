package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// RecaptchaClient is a client for interacting with the reCAPTCHA API
type RecaptchaClient struct {
	secretKey  string
	httpClient *http.Client
}

// RecaptchaResponse represents a response from the reCAPTCHA API
type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

// NewRecaptchaClient creates a new reCAPTCHA client
func NewRecaptchaClient(secretKey string) *RecaptchaClient {
	return &RecaptchaClient{
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// VerifyToken verifies a reCAPTCHA token
func (r *RecaptchaClient) VerifyToken(token string) (bool, error) {
	// If no secret key is set, bypass verification (for development/testing)
	if r.secretKey == "" || r.secretKey == "your-recaptcha-secret-key" {
		return true, nil
	}

	// Make request to reCAPTCHA API
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
		"secret":   {r.secretKey},
		"response": {token},
	})
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Parse response
	var recaptchaResp RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&recaptchaResp); err != nil {
		return false, err
	}

	return recaptchaResp.Success, nil
}
