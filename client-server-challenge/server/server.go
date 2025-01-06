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
}
