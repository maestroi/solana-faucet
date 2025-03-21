package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

// SolanaClient is a client for interacting with the Solana blockchain
type SolanaClient struct {
	rpcClient *rpc.Client
	wallet    solana.PrivateKey
	publicKey solana.PublicKey
}

// NewSolanaClient creates a new Solana client
func NewSolanaClient(rpcURL, walletPath string) (*SolanaClient, error) {
	// Create RPC client
	rpcClient := rpc.New(rpcURL)

	// Load wallet
	data, err := ioutil.ReadFile(walletPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read wallet file: %w", err)
	}

	// Parse wallet data
	var keyBytes []byte
	if err := json.Unmarshal(data, &keyBytes); err != nil {
		return nil, fmt.Errorf("failed to parse wallet: %w", err)
	}

	// Create private key from bytes
	wallet := solana.PrivateKey(keyBytes)
	publicKey := wallet.PublicKey()

	return &SolanaClient{
		rpcClient: rpcClient,
		wallet:    wallet,
		publicKey: publicKey,
	}, nil
}

// GetBalance gets the balance of a wallet in SOL
func (c *SolanaClient) GetBalance(address string) (float64, error) {
	log.Printf("[Solana] Getting balance for address: %s", address)

	// Parse address
	pubKey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		log.Printf("[Solana] Invalid address format: %s", address)
		return 0, fmt.Errorf("invalid Solana address: %w", err)
	}

	// Get balance
	balance, err := c.rpcClient.GetBalance(
		context.Background(),
		pubKey,
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		log.Printf("[Solana] Error getting balance: %v", err)
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	// Convert lamports to SOL
	balanceInSol := float64(balance.Value) / 1e9
	log.Printf("[Solana] Balance for %s: %f SOL (%d lamports)", address, balanceInSol, balance.Value)

	return balanceInSol, nil
}

// IsValidSolanaAddress checks if a string is a valid Solana address
func IsValidSolanaAddress(address string) bool {
	_, err := solana.PublicKeyFromBase58(address)
	return err == nil
}

// SendSOL sends SOL from the faucet wallet to the specified address
func (c *SolanaClient) SendSOL(toAddress string, amount float64) (string, error) {
	log.Printf("[Solana] Sending %f SOL to %s", amount, toAddress)

	// Parse recipient address
	recipient, err := solana.PublicKeyFromBase58(toAddress)
	if err != nil {
		log.Printf("[Solana] Invalid recipient address: %s", toAddress)
		return "", fmt.Errorf("invalid recipient address: %w", err)
	}

	// Convert SOL to lamports
	lamports := uint64(amount * 1e9)

	// Create transfer instruction
	instruction := system.NewTransferInstruction(
		lamports,
		c.publicKey,
		recipient,
	).Build()

	// Get recent blockhash
	recent, err := c.rpcClient.GetLatestBlockhash(context.Background(), rpc.CommitmentConfirmed)
	if err != nil {
		log.Printf("[Solana] Error getting recent blockhash: %v", err)
		return "", fmt.Errorf("failed to get recent blockhash: %w", err)
	}

	// Build transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(c.publicKey),
	)
	if err != nil {
		log.Printf("[Solana] Error creating transaction: %v", err)
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	// Sign transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(c.publicKey) {
			return &c.wallet
		}
		return nil
	})
	if err != nil {
		log.Printf("[Solana] Error signing transaction: %v", err)
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	sig, err := c.rpcClient.SendTransactionWithOpts(
		context.Background(),
		tx,
		rpc.TransactionOpts{
			SkipPreflight:       false,
			PreflightCommitment: rpc.CommitmentConfirmed,
		},
	)
	if err != nil {
		log.Printf("[Solana] Error sending transaction: %v", err)
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	log.Printf("[Solana] Transaction sent: %s", sig.String())
	return sig.String(), nil
}

// GetFaucetBalance returns the balance of the faucet wallet in lamports
func (c *SolanaClient) GetFaucetBalance() (uint64, error) {
	log.Printf("[Solana] Getting faucet wallet balance for address: %s", c.publicKey)

	// Get balance
	balance, err := c.rpcClient.GetBalance(
		context.Background(),
		c.publicKey,
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		log.Printf("[Solana] Error getting balance: %v", err)
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	log.Printf("[Solana] Faucet balance: %d lamports", balance.Value)
	return balance.Value, nil
}
