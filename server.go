package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Struct to store the response from the API
type CotacaoDolarResponse struct {
	USDBRL struct {
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

func NewCotacaoDolarModel(cotacao *CotacaoDolarResponse) *CotacaoDolarModel {
	high, _ := strconv.ParseFloat(cotacao.USDBRL.High, 64)
	low, _ := strconv.ParseFloat(cotacao.USDBRL.Low, 64)
	varBid, _ := strconv.ParseFloat(cotacao.USDBRL.VarBid, 64)
	pctChange, _ := strconv.ParseFloat(cotacao.USDBRL.PctChange, 64)
	bid, _ := strconv.ParseFloat(cotacao.USDBRL.Bid, 64)
	ask, _ := strconv.ParseFloat(cotacao.USDBRL.Ask, 64)
	timestamp, _ := time.Parse(time.RFC3339, cotacao.USDBRL.Timestamp)
	return &CotacaoDolarModel{
		Code:       cotacao.USDBRL.Code,
		Codein:     cotacao.USDBRL.Codein,
		Name:       cotacao.USDBRL.Name,
		High:       high,
		Low:        low,
		VarBid:     varBid,
		PctChange:  pctChange,
		Bid:        bid,
		Ask:        ask,
		Timestamp:  timestamp,
		CreateDate: cotacao.USDBRL.CreateDate,
	}
}

func main() {
	http.HandleFunc("/cotacao", CotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

// Handler to get the cotacao from the API
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
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)

	// Save the response in the database
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&CotacaoDolarModel{})
	var cotacaoModel = NewCotacaoDolarModel(cotacao)

	ctx, cancel = context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancel()

	result := db.WithContext(ctx).Create(cotacaoModel)
	if result.Error != nil {
		fmt.Println("Error saving cotacao in the database")
	}
}

// Function to get the cotacao from the API
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
