package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"net/http"
	"rinhabackend/external/database"
	"rinhabackend/external/lib"
	"rinhabackend/model"
	"rinhabackend/shared/dto"
	"rinhabackend/web/middleware"
	"strings"
	"time"
)

var (
	logger       *slog.Logger
	dbPql        *sql.DB
	limitAccount map[string]int64
)

func init() {
	dbPql = database.ConnPSQL()
	if err := dbPql.Ping(); err != nil {
		log.Fatal(dbPql.Ping())
	}

	accounts, err := getDebitLimitsAccount()
	if err != nil {
		panic(err)
	}

	limitAccount = accounts

	logger = lib.NewLogger(slog.LevelInfo)
	logger.Info("Initialize server")
	logger.Info("started cache account", "accounts", limitAccount)
}

func main() {
	//defer dbPql.Close()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	//r.Use(middleware.SetResponseJsonHeaders)
	//r.Use(middleware.TimeoutLimit)
	r.Get("/", middleware.ProcessTimeout(func(writer http.ResponseWriter, request *http.Request) {
		deadLine, ok := request.Context().Deadline()
		fmt.Println("ContextTimout", deadLine, ok)
		time.Sleep(time.Second * 10)
		fmt.Println("processou")
		writer.Write([]byte("OK"))
	}, 5*time.Second))
	r.Get("/clientes/{id}/extrato", GetTransactions)
	r.Post("/clientes/{id}/transacoes", CreateTransaction)

	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":3000", r))
	// s := &http.Server{Addr: ":3000", Handler: r}
	// log.Fatal(s.ListenAndServeTLS("server.crt", "server.key"))
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "id")

	const queryGetBalancer = `
	SELECT transactions.id, transactions.client_id, transactions.type, transactions.description, transactions.value, transactions.created_at, clients.balance, clients."limit" 
	FROM clients
	LEFT JOIN transactions on clients.id = transactions.client_id
	WHERE clients.id = $1 ORDER BY transactions.created_at DESC LIMIT 10;
	`
	rows, err := dbPql.Query(queryGetBalancer, accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var output dto.GetBalanceOutputDTO
	for rows.Next() {
		var r model.ExtractTransactionModel
		if err := rows.Scan(
			&r.TransactionID,
			&r.ClientID,
			&r.TransactionType,
			&r.TransactionDesc,
			&r.TransactionValue,
			&r.CreatedAt,
			&r.AccountBalance,
			&r.AccountLimit,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		output.Ballance = dto.Ballance{
			Total: int(r.AccountBalance.Int64),
			Limit: int(r.AccountLimit.Int64),
			Date:  time.Now(),
		}
		if r.TransactionValue.Int64 > int64(0) {
			output.Transactions = append(output.Transactions, dto.Transaction{
				Value:     r.TransactionValue.Int64,
				Type:      r.TransactionType.String,
				Desc:      r.TransactionDesc.String,
				CreatedAt: r.CreatedAt.Time,
			})
		} else {
			output.Transactions = []dto.Transaction{}
		}
	}

	if output.Ballance.Total == 0 && output.Ballance.Limit == 0 {
		http.Error(w, "Not found account", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	accountIDParam := chi.URLParam(r, "id")
	defer r.Body.Close()

	var (
		input  dto.Transaction
		output dto.TransactionOutput
	)

	if limitAccount[accountIDParam] == 0 {
		http.Error(w, "Not found account", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Value <= 0 {
		http.Error(w, "Invalid value", http.StatusUnprocessableEntity)
		return
	}
	if input.Type != "c" && input.Type != "d" {
		http.Error(w, "Invalid type", http.StatusUnprocessableEntity)
		return
	}
	if len(input.Desc) == 0 || len(input.Desc) > 10 {
		http.Error(w, "Invalid Description", http.StatusUnprocessableEntity)
		return
	}
	if input.Value > limitAccount[accountIDParam] && input.Type == "d" {
		logger.Warn("Debit: Value must be not greater than the limit",
			"account_id", accountIDParam,
			"cache_limit", limitAccount[accountIDParam],
			"input", input,
		)
		http.Error(w, "Unauthorized transaction E0001x", http.StatusUnprocessableEntity)
		return
	}

	const query = `SELECT "balance", "limit" FROM insert_transactions($1,$2,$3,$4);`
	if err := dbPql.QueryRow(query, accountIDParam, input.Value, input.Type, input.Desc).Scan(
		&output.Balance,
		&output.Limit,
	); err != nil {
		if strings.Contains(err.Error(), "Debit must be positive") {
			http.Error(w, "Unauthorized transaction E0002x", http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func getDebitLimitsAccount() (map[string]int64, error) {
	output := make(map[string]int64)
	const query = `SELECT id, "limit" FROM clients limit 10;`
	rows, err := dbPql.Query(query)
	if err != nil {
		return make(map[string]int64), err
	}
	for rows.Next() {
		var (
			id    string
			limit int64
		)
		if err = rows.Scan(
			&id,
			&limit,
		); err != nil {
			return make(map[string]int64), err
		}
		output[id] = limit
	}
	return output, err
}
