package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type CurrentDollar struct {
	ValueUSD string `json:"bid"`
}

func main() {
	currentDollarValue, err := getCurrentDollarValue()
	if err != nil {
		// handle error
		return
	}
	// use currentDollarValue
	log.Println(currentDollarValue)
}

func getCurrentDollarValue() (*CurrentDollar, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	req, error := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if error != nil {
		return nil, error
	}
	res, error := http.DefaultClient.Do(req)
	if error != nil {
		return nil, error
	}
	defer res.Body.Close()
	var currentDollar CurrentDollar
	error = json.NewDecoder(res.Body).Decode(&currentDollar)
	if error != nil {
		return nil, error
	}
	return &currentDollar, nil
}
