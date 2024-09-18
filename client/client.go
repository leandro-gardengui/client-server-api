package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type CurrentDollar struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	currentDollarValue, err := getCurrentDollarValue()
	if err != nil {
		log.Println(err.Error())
		return
	}
	error := saveResultInFile(currentDollarValue)
	if error != nil {
		log.Println(error.Error())
		return
	}
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

func saveResultInFile(currentDollar *CurrentDollar) error {
	file, error := os.Create("cotacao.txt")
	if error != nil {
		return error
	}
	defer file.Close()
	_, error = file.WriteString("Dólar: " + currentDollar.USDBRL.Bid)
	if error != nil {
		return error
	}
	return nil
}
