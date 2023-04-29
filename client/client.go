package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

type DollarQuote struct {
	Bid string `json:"bid"`
}

func main() {
	dollarQuote, error := requestCotacao()
	if error != nil {
		panic(error)
	}

	error = writeQuotationFile(dollarQuote)
	if error != nil {
		panic(error)
	}
}

func requestCotacao() (*DollarQuote, error) {
	context, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	url := "http://localhost:8080/cotacao"
	request, error := http.NewRequestWithContext(context, "GET", url, nil)
	if error != nil {
		return nil, error
	}

	response, error := http.DefaultClient.Do(request)
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()

	jsonData, error := io.ReadAll(response.Body)
	if error != nil {
		return nil, error
	}
	var dollarQuote DollarQuote
	error = json.Unmarshal(jsonData, &dollarQuote)
	if error != nil {
		return nil, error
	}

	return &dollarQuote, nil
}

func writeQuotationFile(dollarQuote *DollarQuote) error {
	file, error := os.Create("cotacao.txt")
	if error != nil {
		return error
	}
	_, err := file.Write([]byte("Dólar:" + dollarQuote.Bid + " \n"))
	if err != nil {
		return err
	}

	return nil
}
