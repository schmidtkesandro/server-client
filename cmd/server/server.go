package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

type USDBRL struct {
	Currency Currency `json:"USDBRL"`
}

func main() {
	db, err := sql.Open("sqlite3", "./database/cotacoes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// cria tabela para armazenar as cotações
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (
        id INTEGER PRIMARY KEY,
        datahora DATETIME DEFAULT CURRENT_TIMESTAMP,
        valor FLOAT
    );`)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		//define timeout para a solicitação da API de cotação

		//Crie um contexto com um tempo limite para a solicitação HTTP
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		//Crie uma nova solicitação HTTP com o contexto e o URL da API
		req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
		if err != nil {
			log.Fatal("err", err)
		}

		// Crie um cliente HTTP padrão e envie a solicitação HTTP
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal("erro", err)
		}
		defer resp.Body.Close()

		// =================================================================
		// Decodifique a resposta JSON em uma struct
		var data map[string]json.RawMessage
		var currency Currency

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Fatal("erro0", err)
		}

		// Acesse o valor desejado pelo índice da estrutura
		usdbrl := data["USDBRL"]

		err = json.Unmarshal(usdbrl, &currency)
		if err != nil {
			log.Fatal("erro1", err)
		}
		//================================================================
		fmt.Println(string(usdbrl))
		fmt.Println("Valor do dólar : ", currency.Bid)
		fmt.Println("Valor do code : ", currency.Code)
		fmt.Println("Valor do name : ", currency.Name)
		ctx, cancel = context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancel()

		// Cria a instrução SQL para inserir a cotação no banco de dados
		stmt, err := db.PrepareContext(ctx, "INSERT INTO cotacoes (valor) VALUES (?)")
		if err != nil {
			log.Fatal("erro BD\n", err)
			return
		}
		defer stmt.Close()

		// Executa a instrução SQL
		_, err = stmt.ExecContext(ctx, currency.Bid)
		if err != nil {
			log.Fatal("erro BD\n", err)
			return
		}

		// retorna o valor da cotação em formato JSON para o cliente
		w.Header().Set("Content-Type", "application/json")
		//json.NewEncoder(w).Encode(usdbrl)
		json.NewEncoder(w).Encode(currency)
	})

	fmt.Println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}
