package db

import (
	"database/sql"
	"time"

	"github.com/maestroi/solana-faucet/backend/models"
	_ "github.com/mattn/go-sqlite3"
)

// Database represents the application database
type Database struct {
	db *sql.DB
}

// migrateTimestamps converts existing timestamps to ISO format
func (d *Database) migrateTimestamps() error {
	// First, check if we have any timestamps in the old format
	query := `
	SELECT id, timestamp
	FROM transactions
	WHERE timestamp NOT LIKE '%Z' AND timestamp NOT LIKE '%+%'
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Prepare update statement
	updateStmt, err := d.db.Prepare(`
		UPDATE transactions
		SET timestamp = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer updateStmt.Close()

	// Process each row
	for rows.Next() {
		var id int64
		var timestamp string
		if err := rows.Scan(&id, &timestamp); err != nil {
			return err
		}

		// Parse the old timestamp
		t, err := time.Parse("2006-01-02 15:04:05", timestamp)
		if err != nil {
			continue // Skip if we can't parse
		}

		// Convert to ISO format
		isoTimestamp := t.Format(time.RFC3339)

		// Update the row
		if _, err := updateStmt.Exec(isoTimestamp, id); err != nil {
			return err
		}
	}

	return nil
}

// migrateClaimHistoryTimestamps converts existing timestamps to ISO format
func (d *Database) migrateClaimHistoryTimestamps() error {
	// First, check if we have any timestamps in the old format
	query := `
	SELECT id, last_claim_time
	FROM claim_history
	WHERE last_claim_time NOT LIKE '%Z' AND last_claim_time NOT LIKE '%+%'
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Prepare update statement
	updateStmt, err := d.db.Prepare(`
		UPDATE claim_history
		SET last_claim_time = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer updateStmt.Close()

	// Process each row
	for rows.Next() {
		var id int64
		var lastClaimTime string
		if err := rows.Scan(&id, &lastClaimTime); err != nil {
			return err
		}

		// Parse the old timestamp
		t, err := time.Parse("2006-01-02 15:04:05", lastClaimTime)
		if err != nil {
			continue // Skip if we can't parse
		}

		// Convert to ISO format
		isoTimestamp := t.Format(time.RFC3339)

		// Update the row
		if _, err := updateStmt.Exec(isoTimestamp, id); err != nil {
			return err
		}
	}

	return nil
}

// InitDB initializes the database
func InitDB(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, err
	}

	database := &Database{db: db}

	// Migrate timestamps to ISO format
	if err := database.migrateTimestamps(); err != nil {
		return nil, err
	}

	// Migrate claim history timestamps to ISO format
	if err := database.migrateClaimHistoryTimestamps(); err != nil {
		return nil, err
	}

	return database, nil
}

// createTables creates the necessary tables in the database
func createTables(db *sql.DB) error {
	// Create transactions table
	transactionTableSQL := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		wallet_address TEXT NOT NULL,
		ip_address TEXT NOT NULL,
		amount REAL NOT NULL,
		status TEXT NOT NULL,
		tx_hash TEXT,
		error_message TEXT,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(transactionTableSQL); err != nil {
		return err
	}

	// Create claim_history table
	claimHistoryTableSQL := `
	CREATE TABLE IF NOT EXISTS claim_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		wallet_address TEXT NOT NULL,
		ip_address TEXT NOT NULL,
		last_claim_time TIMESTAMP NOT NULL,
		claim_count INTEGER NOT NULL DEFAULT 1,
		UNIQUE(wallet_address)
	);
	`
	if _, err := db.Exec(claimHistoryTableSQL); err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// GetClaimHistory retrieves the claim history for a wallet
func (d *Database) GetClaimHistory(walletAddress string) (*models.ClaimHistory, error) {
	query := `
	SELECT id, wallet_address, ip_address, last_claim_time, claim_count
	FROM claim_history
	WHERE wallet_address = ?
	`

	row := d.db.QueryRow(query, walletAddress)

	var ch models.ClaimHistory
	var lastClaimTime string

	err := row.Scan(&ch.ID, &ch.WalletAddress, &ch.IPAddress, &lastClaimTime, &ch.ClaimCount)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse the timestamp using RFC3339 format
	t, err := time.Parse(time.RFC3339, lastClaimTime)
	if err != nil {
		// Try the old format as fallback
		t, err = time.Parse("2006-01-02 15:04:05", lastClaimTime)
		if err != nil {
			return nil, err
		}
	}
	ch.LastClaimTime = t

	return &ch, nil
}

// GetClaimHistoryByIP retrieves the claim history for an IP address
func (d *Database) GetClaimHistoryByIP(ipAddress string) ([]*models.ClaimHistory, error) {
	query := `
	SELECT id, wallet_address, ip_address, last_claim_time, claim_count
	FROM claim_history
	WHERE ip_address = ?
	`

	rows, err := d.db.Query(query, ipAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []*models.ClaimHistory

	for rows.Next() {
		var ch models.ClaimHistory
		var lastClaimTime string

		if err := rows.Scan(&ch.ID, &ch.WalletAddress, &ch.IPAddress, &lastClaimTime, &ch.ClaimCount); err != nil {
			return nil, err
		}

		// Parse the timestamp using RFC3339 format
		t, err := time.Parse(time.RFC3339, lastClaimTime)
		if err != nil {
			// Try the old format as fallback
			t, err = time.Parse("2006-01-02 15:04:05", lastClaimTime)
			if err != nil {
				return nil, err
			}
		}
		ch.LastClaimTime = t

		histories = append(histories, &ch)
	}

	return histories, nil
}

// UpdateClaimHistory updates or creates a claim history record
func (d *Database) UpdateClaimHistory(walletAddress, ipAddress string) error {
	// Check if record exists
	history, err := d.GetClaimHistory(walletAddress)
	if err != nil {
		return err
	}

	if history == nil {
		// Create new record
		query := `
		INSERT INTO claim_history (wallet_address, ip_address, last_claim_time, claim_count)
		VALUES (?, ?, strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), 1)
		`
		_, err := d.db.Exec(query, walletAddress, ipAddress)
		return err
	}

	// Update existing record
	query := `
	UPDATE claim_history
	SET last_claim_time = strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), 
		ip_address = ?,
		claim_count = claim_count + 1
	WHERE wallet_address = ?
	`
	_, err = d.db.Exec(query, ipAddress, walletAddress)
	return err
}

// CreateTransaction creates a new transaction record
func (d *Database) CreateTransaction(tx *models.Transaction) (int64, error) {
	query := `
	INSERT INTO transactions (wallet_address, ip_address, amount, status, tx_hash, error_message, timestamp)
	VALUES (?, ?, ?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
	`

	result, err := d.db.Exec(
		query,
		tx.WalletAddress,
		tx.IPAddress,
		tx.Amount,
		tx.Status,
		tx.TxHash,
		tx.ErrorMessage,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// UpdateTransaction updates an existing transaction
func (d *Database) UpdateTransaction(tx *models.Transaction) error {
	query := `
	UPDATE transactions
	SET status = ?, tx_hash = ?, error_message = ?
	WHERE id = ?
	`

	_, err := d.db.Exec(query, tx.Status, tx.TxHash, tx.ErrorMessage, tx.ID)
	return err
}

// GetRecentTransactions retrieves recent transactions
func (d *Database) GetRecentTransactions(limit int) ([]*models.Transaction, error) {
	query := `
	SELECT id, wallet_address, amount, status, tx_hash, error_message, timestamp
	FROM transactions
	ORDER BY timestamp DESC
	LIMIT ?
	`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction

	for rows.Next() {
		var tx models.Transaction
		var timestamp string

		if err := rows.Scan(
			&tx.ID,
			&tx.WalletAddress,
			&tx.Amount,
			&tx.Status,
			&tx.TxHash,
			&tx.ErrorMessage,
			&timestamp,
		); err != nil {
			return nil, err
		}

		// Parse the timestamp using RFC3339 format
		t, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			// Try the old format as fallback
			t, err = time.Parse("2006-01-02 15:04:05", timestamp)
			if err != nil {
				return nil, err
			}
		}
		tx.Timestamp = t

		transactions = append(transactions, &tx)
	}

	return transactions, nil
}

// GetTransactionsByWallet retrieves transactions for a specific wallet
func (d *Database) GetTransactionsByWallet(walletAddress string, limit int) ([]*models.Transaction, error) {
	query := `
	SELECT id, wallet_address, amount, status, tx_hash, error_message, timestamp
	FROM transactions
	WHERE wallet_address = ?
	ORDER BY timestamp DESC
	LIMIT ?
	`

	rows, err := d.db.Query(query, walletAddress, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction

	for rows.Next() {
		var tx models.Transaction
		var timestamp string

		if err := rows.Scan(
			&tx.ID,
			&tx.WalletAddress,
			&tx.Amount,
			&tx.Status,
			&tx.TxHash,
			&tx.ErrorMessage,
			&timestamp,
		); err != nil {
			return nil, err
		}

		// Parse the timestamp using RFC3339 format
		t, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			// Try the old format as fallback
			t, err = time.Parse("2006-01-02 15:04:05", timestamp)
			if err != nil {
				return nil, err
			}
		}
		tx.Timestamp = t

		transactions = append(transactions, &tx)
	}

	return transactions, nil
}
