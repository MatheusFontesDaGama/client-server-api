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
	dollarQuote, errorRequestCotacao := requestCotacao()
	if errorRequestCotacao != nil {
		panic(errorRequestCotacao)
	}

	errorWriteFile := writeQuotationFile(dollarQuote)
	if errorWriteFile != nil {
		panic(errorWriteFile)
	}
}

func requestCotacao() (*DollarQuote, error) {
	context, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	url := "http://localhost:8080/cotacao"
	request, errorRequest := http.NewRequestWithContext(context, "GET", url, nil)
	if errorRequest != nil {
		return nil, errorRequest
	}

	response, errorResponse := http.DefaultClient.Do(request)
	if errorResponse != nil {
		return nil, errorResponse
	}
	defer response.Body.Close()

	jsonData, errorReadAll := io.ReadAll(response.Body)
	if errorReadAll != nil {
		return nil, errorReadAll
	}
	var dollarQuote DollarQuote
	errorUnmarshal := json.Unmarshal(jsonData, &dollarQuote)
	if errorUnmarshal != nil {
		return nil, errorUnmarshal
	}

	return &dollarQuote, nil
}

func writeQuotationFile(dollarQuote *DollarQuote) error {
	file, errorCreateFile := os.Create("../cotacao.txt")
	if errorCreateFile != nil {
		return errorCreateFile
	}
	_, errorWriteFile := file.Write([]byte("Dólar:" + dollarQuote.Bid + " \n"))
	if errorWriteFile != nil {
		return errorWriteFile
	}

	return nil
}
