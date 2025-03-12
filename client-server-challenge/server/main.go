package main

import (
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

	// Fetch the currency data from the API
	response, err := http.Get(apiURL)
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

	// Parse the response data
	var currencyData CurrencyData
	if err = json.Unmarshal(responseData, &currencyData); err != nil {
		handleError(w, err, http.StatusInternalServerError, "Failed to parse JSON")
		return
	}

	// Insert the currency data into the database
	if err = insertData(currencyData); err != nil {
		handleError(w, err, http.StatusInternalServerError, "Failed to store currency data")
		return
	}

	// Write the JSON response to the client
	writeJSONResponse(w, currencyData)
}

// Insert data function to insert the currency data into the database
func insertData(currencyData CurrencyData) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}

	// Insert the currency data into the database
	_, err := db.Exec(
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

	// Handle the error if the data is not inserted
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	return nil
}

// Get data handler function to get the currency data from the database
func getDataHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the error if the method is not GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle the error if the database is not initialized
	if db == nil {
		handleError(w, errors.New("database not initialized"), http.StatusInternalServerError, "Database connection is not initialized")
		return
	}

	// Execute the query to get the currency data from the database
	rows, err := db.Query(getAllDataSQL)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Failed to execute query")
		return
	}
	defer rows.Close()

	// Create a slice to store the currency data
	var results []ExchangeRate
	for rows.Next() {
		// Create a variable to store the currency data
		var result ExchangeRate
		err = rows.Scan(
			&result.Code,
			&result.Codein,
			&result.Name,
			&result.High,
			&result.Low,
			&result.VarBid,
			&result.PctChange,
			&result.Bid,
			&result.Ask,
			&result.Timestamp,
			&result.CreateDate,
		)
		// Handle the error if the data is not scanned
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "Failed to scan row")
			return
		}
		results = append(results, result)
	}

	// Handle the error if the data is not iterated over
	if err = rows.Err(); err != nil {
		handleError(w, err, http.StatusInternalServerError, "Failed to iterate over rows")
		return
	}

	// Write the JSON response to the client
	writeJSONResponse(w, results)
}

// Setup routes function to setup the routes for the server
func setupRoutes() {
	http.HandleFunc("/cotacao", fetchCurrencyHandler)
	http.HandleFunc("/get-data", getDataHandler)
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
