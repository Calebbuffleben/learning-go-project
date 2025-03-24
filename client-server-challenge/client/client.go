package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

	fmt.Printf("Fetching data from server at %s...\n", serverURL)

	err := fetchCurrencyData(serverURL + "/cotacao") // Fetch currency data
	if err != nil {                                  // If there's an error
		log.Printf("Failed to fetch currency data: %v", err) // Log the error
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

	// Log raw response for debugging
	fmt.Printf("Raw response: %s\n", string(body))

	// The server returns just a simple JSON with a bid field
	var bidResponse struct {
		Bid string `json:"bid"`
	}
	err = json.Unmarshal(body, &bidResponse) // Unmarshal the JSON data
	if err != nil {                          // If there's an error
		return fmt.Errorf("failed to unmarshal JSON: %w", err) // Return the error
	}

	fmt.Printf("Fetched bid rate: %s\n", bidResponse.Bid) // Print the bid rate
	
	// Save the bid to a file
	err = saveBidToFile(bidResponse.Bid)
	if err != nil {
		return fmt.Errorf("failed to save bid to file: %w", err)
	}
	
	return nil // Return nil error
}

// Function to save bid rate to a file
func saveBidToFile(bid string) error {
	// Create a file to save the bid rate
	file, err := os.Create("./cotacao.txt")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write the bid rate to the file
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s\n", bid))
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Println("Bid rate saved to cotacao.txt")
	return nil
}