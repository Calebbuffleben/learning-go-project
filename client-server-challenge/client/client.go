package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Define a struct to hold exchange rate data
type ExchangeRate struct {
	Code       string `json:"code"`        // Currency code
	Codein     string `json:"codein"`      // Base currency code
	Name       string `json:"name"`        // Currency name
	High       string `json:"high"`        // Highest price
	Low        string `json:"low"`         // Lowest price
	VarBid     string `json:"varBid"`      // Variation of the bid price
	PctChange  string `json:"pctChange"`   // Percentage change
	Bid        string `json:"bid"`         // Bid price
	Ask        string `json:"ask"`         // Ask price
	Timestamp  string `json:"timestamp"`   // Timestamp of the data
	CreateDate string `json:"create_date"` // Date the data was created
}

// Main function
func main() {
	// Get server URL from environment variable, with a default fallback
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080"
	}

	for { // Infinite loop
		fmt.Printf("Fetching data from server at %s...\n", serverURL)

		err := fetchCurrencyData(serverURL + "/cotacao") // Fetch currency data
		if err != nil {                                  // If there's an error
			log.Printf("Failed to fetch currency data: %v", err) // Log the error
		}

		fmt.Println("Fetching stored data...")         // Print a message
		err = fetchStoredData(serverURL + "/get-data") // Fetch stored data
		if err != nil {                                // If there's an error
			log.Printf("Failed to fetch stored data: %v", err) // Log the error
		}

		// Add a small delay to prevent overwhelming the server
		time.Sleep(10 * time.Second)
	}
}

// Function to fetch currency data
func fetchCurrencyData(url string) error {
	resp, err := http.Get(url) // Make a GET request to the URL
	if err != nil {            // If there's an error
		return fmt.Errorf("failed to fetch data: %w", err) // Return the error
	}
	defer resp.Body.Close() // Close the response body when the function returns

	body, err := io.ReadAll(resp.Body) // Read the response body
	if err != nil {                    // If there's an error
		return fmt.Errorf("failed to read response body: %w", err) // Return the error
	}

	var data map[string]ExchangeRate  // Declare a map to hold the data
	err = json.Unmarshal(body, &data) // Unmarshal the JSON data into the map
	if err != nil {                   // If there's an error
		return fmt.Errorf("failed to unmarshal JSON: %w", err) // Return the error
	}

	fmt.Printf("Fetched currency data: %+v\n", data) // Print the fetched data
	return nil                                       // Return nil error
}

// Function to fetch stored data
func fetchStoredData(url string) error {
	resp, err := http.Get(url) // Make a GET request to the URL
	if err != nil {            // If there's an error
		return fmt.Errorf("failed to fetch data: %w", err) // Return the error
	}
	defer resp.Body.Close() // Close the response body when the function returns

	body, err := io.ReadAll(resp.Body) // Read the response body
	if err != nil {                    // If there's an error
		return fmt.Errorf("failed to read response body: %w", err) // Return the error
	}

	var data []ExchangeRate           // Declare a slice to hold the data
	err = json.Unmarshal(body, &data) // Unmarshal the JSON data into the slice
	if err != nil {                   // If there's an error
		return fmt.Errorf("failed to unmarshal JSON: %w", err) // Return the error
	}

	fmt.Printf("Fetched stored data: %+v\n", data) // Print the fetched data
	return nil                                     // Return nil error
}
