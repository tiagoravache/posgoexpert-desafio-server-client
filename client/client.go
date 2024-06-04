package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	serverAddress = "http://localhost:8080"
	clientTimeout = 300 * time.Millisecond
	dataFolder    = "data/"
	fileName      = "quotation.txt"
)

type Quotation struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	quotation, err := fetchQuotation(ctx)
	if err != nil {
		log.Fatal("Erro ao obter cotação: ", err)
	}

	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		os.Mkdir(dataFolder, 0755)
	}

	if err := saveToFile(quotation, dataFolder+fileName); err != nil {
		log.Fatal("Erro salvando cotação: ", err)
	}
}

func fetchQuotation(ctx context.Context) (string, error) {
	client := &http.Client{}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/cotacao", serverAddress), nil)
	if err != nil {
		return "", err
	}

	// Do request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Código de status inesperado: %d", resp.StatusCode)
	}

	// Decode response
	var response Quotation
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Print(err)
		return "", err
	}

	return response.Bid, nil
}

func saveToFile(bid, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	timestamp := time.Now().Format(time.DateTime)
	_, err = fmt.Fprintf(file, "%s - Dólar: %s\n", timestamp, bid)
	return err
}
