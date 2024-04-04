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

type Currency struct {
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
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

func main() {
	// faz a solicitação da cotação para o servidor
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// Decodifique a resposta JSON em uma struct
	//var data map[string]json.RawMessage
	var currency Currency
	if err := json.NewDecoder(resp.Body).Decode(&currency); err != nil {
		log.Fatal(err)
	}

	// Abre o arquivo "cotacao.txt" para escrita com adição (append) ou criação se não existir
	file, err := os.OpenFile("./files/cotacao.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	// Monta a mensagem que será gravada no arquivo
	content := "\nDólar:" + currency.Bid + " cotado em " + time.Now().Format("2006-01-02 15:04:05")

	// Escreve a string formatada no arquivo
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo:", err)
		return
	}
	// Mostra na linha de comando o que será gravado no arquivo
	fmt.Println(content)
}
