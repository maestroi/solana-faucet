package models

// FundRequest represents a request for funds
type FundRequest struct {
	WalletAddress     string `json:"wallet_address"`
	TurnstileResponse string `json:"cf_turnstile_response"`
}
