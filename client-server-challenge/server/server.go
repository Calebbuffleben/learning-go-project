package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
}

func handleRequests() {
	http.Handle("/getCoinValue", http.HandlerFunc(server))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func main() {
	handleRequests()

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

	// Sample data
	data := CurrencyData{
		Code:       "USD",
		CodeIn:     "BRL",
		Name:       "Dólar Americano/Real Brasileiro",
		High:       6.2012,
		Low:        6.0985,
		VarBid:     -0.0695,
		PctChange:  -1.12,
		Bid:        6.1098,
		Ask:        6.1108,
		Timestamp:  1736197196,
		CreateDate: "2025-01-06 17:59:56",
	}
}
