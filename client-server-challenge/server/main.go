package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Constants
const (
	dbPath = "./db/database.db"
	apiURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
)

// CurrencyData struct to store the currency data
type CurrencyData struct {
	USDBRL ExchangeRate `json:"USDBRL"`
}

// ExchangeRate struct to store the currency data
type ExchangeRate struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

// Package-level variable for the database connection
var db *sql.DB

// SQL queries
const (
	createTableSQL = `
	CREATE TABLE IF NOT EXISTS currency_data (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT NOT NULL,
		codein TEXT NOT NULL,
		name TEXT NOT NULL,
		high REAL NOT NULL,
		low REAL NOT NULL,
		varBid REAL NOT NULL,
		pctChange REAL NOT NULL,
		bid REAL NOT NULL,
		ask REAL NOT NULL,
		timestamp INTEGER NOT NULL,
		create_date TEXT NOT NULL
	);`

	insertDataSQL = `
		INSERT INTO currency_data (
			code, codein, name, high, low, varBid, pctChange,
			bid, ask, timestamp, create_date
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	getAllDataSQL = `
		SELECT code, codein, name, high, low, varBid, pctChange,
		       bid, ask, timestamp, create_date 
		FROM currency_data;`
)

// Initialization function
func init() {
	if err := setupDatabase(); err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}
}

// Setup database function
func setupDatabase() error {
	if err := ensureDirectory("./db"); err != nil {
		return fmt.Errorf("failed to ensure directory: %w", err)
	}

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create table if it doesn't exist
	if _, err = db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// Ensure directory function to create the db directory if it doesn't exist
func ensureDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

// Write JSON response function to write the JSON response to the client
func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

// Handle error function to handle the error and return the error message to the client
func handleError(w http.ResponseWriter, err error, status int, message string) {
	log.Printf("%s: %v", message, err)
	http.Error(w, message, status)
}

// Fetch currency handler function to fetch the currency data from the API and store it in the database
func fetchCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create a context with timeout for API operations - 200ms timeout
	apiCtx, apiCancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer apiCancel()

	// Create a request with the API context
	req, err := http.NewRequestWithContext(apiCtx, http.MethodGet, apiURL, nil)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Failed to create request")
		return
	}

	// Fetch the currency data from the API with timeout context
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Failed to fetch currency data")
		return
	}
	defer response.Body.Close()

	// Read the response data
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Failed to read response data")
		return
	}

	// Log the raw response for debugging
	log.Printf("Raw API response: %s", string(responseData))

	// First try to parse as a generic map to understand the structure
	var rawMap map[string]interface{}
	if err = json.Unmarshal(responseData, &rawMap); err != nil {
		handleError(w, err, http.StatusInternalServerError, "Failed to parse raw JSON")
		return
	}

	// Create a valid CurrencyData struct
	var currencyData CurrencyData

	// Check if USDBRL exists and what type it is
	if usdbrlVal, exists := rawMap["USDBRL"]; exists {
		switch v := usdbrlVal.(type) {
		case string:
			// USDBRL is a string (likely the bid value)
			currencyData.USDBRL = ExchangeRate{
				Code:       "USD",
				Codein:     "BRL",
				Name:       "Dollar/Real",
				High:       "0",
				Low:        "0",
				VarBid:     "0",
				PctChange:  "0",
				Bid:        v,
				Ask:        "0",
				Timestamp:  fmt.Sprintf("%d", time.Now().Unix()),
				CreateDate: time.Now().Format(time.RFC3339),
			}
		case map[string]interface{}:
			// USDBRL is an object, try to extract fields
			if bidVal, ok := v["bid"].(string); ok {
				currencyData.USDBRL = ExchangeRate{
					Code:       getStringOrDefault(v, "code", "USD"),
					Codein:     getStringOrDefault(v, "codein", "BRL"),
					Name:       getStringOrDefault(v, "name", "Dollar/Real"),
					High:       getStringOrDefault(v, "high", "0"),
					Low:        getStringOrDefault(v, "low", "0"),
					VarBid:     getStringOrDefault(v, "varBid", "0"),
					PctChange:  getStringOrDefault(v, "pctChange", "0"),
					Bid:        bidVal,
					Ask:        getStringOrDefault(v, "ask", "0"),
					Timestamp:  getStringOrDefault(v, "timestamp", fmt.Sprintf("%d", time.Now().Unix())),
					CreateDate: getStringOrDefault(v, "create_date", time.Now().Format(time.RFC3339)),
				}
			} else {
				// No bid value found, create default
				currencyData.USDBRL = createDefaultExchangeRate()
			}
		default:
			// Unknown type, create default
			currencyData.USDBRL = createDefaultExchangeRate()
		}
	} else {
		// No USDBRL key, try standard unmarshal
		if err = json.Unmarshal(responseData, &currencyData); err != nil {
			// If that fails too, create default
			currencyData.USDBRL = createDefaultExchangeRate()
		}
	}

	// Validate and ensure we have a bid value
	if currencyData.USDBRL.Bid == "" {
		currencyData.USDBRL.Bid = "0"
	}

	// Create a new context with timeout for database operations - strict 10ms timeout
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer dbCancel()

	// Insert the currency data into the database using the database context
	if err = insertDataWithContext(dbCtx, currencyData); err != nil {
		// If there's a database error, log it but continue to return the bid
		log.Printf("Database insertion error: %v", err)
	}

	// Create a clean response with just the bid value
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Use json.Marshal to ensure proper JSON formatting
	bidResponse, _ := json.Marshal(map[string]string{"bid": currencyData.USDBRL.Bid})
	w.Write(bidResponse)
}

// Helper function to get string value from map or default if not found
func getStringOrDefault(data map[string]interface{}, key, defaultValue string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return defaultValue
}

// Helper function to create default exchange rate
func createDefaultExchangeRate() ExchangeRate {
	return ExchangeRate{
		Code:       "USD",
		Codein:     "BRL",
		Name:       "Dollar/Real",
		High:       "0",
		Low:        "0",
		VarBid:     "0",
		PctChange:  "0",
		Bid:        "0",
		Ask:        "0",
		Timestamp:  fmt.Sprintf("%d", time.Now().Unix()),
		CreateDate: time.Now().Format(time.RFC3339),
	}
}

// Insert data with context function to insert the currency data into the database
func insertDataWithContext(ctx context.Context, currencyData CurrencyData) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	// Create a channel to handle the result of the database insertion
	done := make(chan error, 1)

	go func() {
		// Insert the currency data into the database with context
		_, err := db.ExecContext(ctx,
			insertDataSQL,
			currencyData.USDBRL.Code,
			currencyData.USDBRL.Codein,
			currencyData.USDBRL.Name,
			currencyData.USDBRL.High,
			currencyData.USDBRL.Low,
			currencyData.USDBRL.VarBid,
			currencyData.USDBRL.PctChange,
			currencyData.USDBRL.Bid,
			currencyData.USDBRL.Ask,
			currencyData.USDBRL.Timestamp,
			currencyData.USDBRL.CreateDate,
		)
		done <- err
	}()

	// Wait for either the context to timeout or the insertion to complete
	select {
	case <-ctx.Done():
		return fmt.Errorf("database insertion timed out: %w", ctx.Err())
	case err := <-done:
		if err != nil {
			return fmt.Errorf("failed to insert data: %w", err)
		}
		return nil
	}
}

// Setup routes function to setup the routes for the server
func setupRoutes() {
	http.HandleFunc("/cotacao", fetchCurrencyHandler)
}

// Main function to start the server
func main() {
	// Setup the routes
	setupRoutes()
	// Print the server started message
	fmt.Println("Server started on :8080")
	// Close the database connection
	defer db.Close()
	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
