package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func server() {
	response, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		fmt.Println(err)
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))
}
