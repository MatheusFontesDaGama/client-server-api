package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type USDBRL struct {
	Usdbrl DollarQuote `json:"USDBRL"`
}

type DollarQuote struct {
	Id         string
	Code       string  `json:"code"`
	CodeIn     string  `json:"codein"`
	Name       string  `json:"name"`
	High       float64 `json:"high,string"`
	Low        float64 `json:"low,string"`
	VarBid     float64 `json:"varBid,string"`
	PctChange  float64 `json:"pctChange,string"`
	Bid        float64 `json:"bid,string"`
	Ask        float64 `json:"ask,string"`
	Timestamp  string  `json:"timestamp"`
	CreateDate string  `json:"create_date"`
}

func main() {
	http.HandleFunc("/cotacao", GetDollarQuoteHandler)
	http.ListenAndServe(":8080", nil)
}

func GetDollarQuoteHandler(response http.ResponseWriter, request *http.Request) {
	dollarQuote, errorDollarQuote := getDollarQuote()
	if errorDollarQuote != nil {
		response.WriteHeader(http.StatusInternalServerError)
		panic(errorDollarQuote)
	}

	db, errorDB := sql.Open("sqlite3", "../dollar_quote.db")
	if errorDB != nil {
		response.WriteHeader(http.StatusInternalServerError)
		panic(errorDB)
	}
	defer db.Close()

	errorInsertDollarQuote := insertDollarQuote(db, dollarQuote)
	if errorInsertDollarQuote != nil {
		response.WriteHeader(http.StatusInternalServerError)
		panic(errorInsertDollarQuote)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(dollarQuote)
}

func getDollarQuote() (*DollarQuote, error) {
	context, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	request, errorRequest := http.NewRequestWithContext(context, "GET", url, nil)
	if errorRequest != nil {
		return nil, errorRequest
	}

	var dollarQuote USDBRL
	responseHttp, errorResponse := http.DefaultClient.Do(request)
	if errorResponse != nil {
		return nil, errorResponse
	}
	defer responseHttp.Body.Close()

	response, errorReadAll := io.ReadAll(responseHttp.Body)
	if errorReadAll != nil {
		return nil, errorReadAll
	}

	errorUnmarshal := json.Unmarshal(response, &dollarQuote)
	if errorUnmarshal != nil {
		return nil, errorUnmarshal
	}

	dollarQuote.Usdbrl.Id = uuid.New().String()
	return &dollarQuote.Usdbrl, nil
}

func insertDollarQuote(db *sql.DB, dollarQuote *DollarQuote) error {
	contextDB, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	sql := "INSERT INTO dollar_quotes(id, code, code_in, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, errorStmt := db.PrepareContext(contextDB, sql)
	if errorStmt != nil {
		return errorStmt
	}
	defer stmt.Close()

	_, errorResult := stmt.ExecContext(
		contextDB,
		dollarQuote.Id,
		dollarQuote.Code,
		dollarQuote.CodeIn,
		dollarQuote.Name,
		dollarQuote.High,
		dollarQuote.Low,
		dollarQuote.VarBid,
		dollarQuote.PctChange,
		dollarQuote.Bid,
		dollarQuote.Ask,
		dollarQuote.Timestamp,
		dollarQuote.CreateDate,
	)
	if errorResult != nil {
		return errorResult
	}

	return nil
}
