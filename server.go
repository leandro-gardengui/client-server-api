package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// Struct to store the response from the API
type CotacaoDolarResponse struct {
	Usdbrl struct {
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
	} `json:"USDBRL"`
}

// Struct to store the model of the database
type CotacaoDolarModel struct {
	ID         int       `gorm:"primaryKey"`
	Code       string    `json:"code"`
	Codein     string    `json:"codein"`
	Name       string    `json:"name"`
	High       float64   `json:"high"`
	Low        float64   `json:"low"`
	VarBid     float64   `json:"varBid"`
	PctChange  float64   `json:"pctChange"`
	Bid        float64   `json:"bid"`
	Ask        float64   `json:"ask"`
	Timestamp  time.Time `json:"timestamp"`
	CreateDate string    `json:"create_date"`
	gorm.Model
}

func main() {
	http.HandleFunc("/cotacao", CotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()
	cotacao, error := getCotacao(ctx)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)
}

func getCotacao(ctx context.Context) (*CotacaoDolarResponse, error) {
	resp, error := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	select {
	default:
		if error != nil {
			return nil, error
		}
		defer resp.Body.Close()
		body, error := io.ReadAll(resp.Body)
		if error != nil {
			return nil, error
		}
		var cotacao CotacaoDolarResponse
		error = json.Unmarshal(body, &cotacao)
		if error != nil {
			return nil, error
		}
		return &cotacao, nil
	case <-ctx.Done():
		fmt.Println("Chamada na API de cotação cancelada. Timeout reached.")
		return nil, ctx.Err()
	}
}
