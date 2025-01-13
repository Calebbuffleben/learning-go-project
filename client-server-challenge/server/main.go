package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"database/sql"
	"os"	

	_ "github.com/mattn/go-sqlite3"
)

type CurrencyData struct {
	USDBRL ExchangeRate `json:"USDBRL"`
}

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

func server(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var currencyData CurrencyData
	err = json.Unmarshal(responseData, &currencyData)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currencyData)
	insertData(currencyData)
}

func insertData(currencyData CurrencyData) error {
	insertDataQuery := `
		INSERT INTO currency_data (
			code,
			codein,
			name,
			high,
			low,
			varBid,
			pctChange,
			bid,
			ask,
			timestamp,
			create_date
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := db.Exec(
		insertDataQuery,
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
	if err != nil {
		return err
	}

	return nil
}
func handleRequests() {
	http.Handle("/cotacao", http.HandlerFunc(server))
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	fmt.Println("Server started")

	// Database file path
	dbPath := "./db/database.db"

	// Ensure the db directory exists
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		err = os.Mkdir("./db", os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create db directory: %v", err)
		}
	}

	// Open SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	createTableQuery := `
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
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	handleRequests()
}
