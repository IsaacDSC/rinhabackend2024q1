package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"log/slog"
	"net/http"
	"rinhabackend/external/database"
	"rinhabackend/external/lib"
	"rinhabackend/external/queue"
	"rinhabackend/shared/dto"
	"rinhabackend/web/middleware"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

var (
	dbCache *Database
	logger  *slog.Logger
	dbPql   *sql.DB
	iSyncQ  *queue.Event
)

func init() {
	dbCache = NewDatabase()
	dbCache.Start()

	dbPql = database.ConnPSQL()
	if err := dbPql.Ping(); err != nil {
		log.Fatal(dbPql.Ping())
	}

	logger = lib.NewLogger(slog.LevelError)
	logger.Info("Initialize server")
}

type Database struct {
	mu    sync.Mutex
	Table map[string]dto.Account
}

func NewDatabase() *Database {
	return new(Database)
}
func (c *Database) Start() {
	c.Table = map[string]dto.Account{
		"1": {
			ID:           1,
			Balance:      0,
			Limit:        100000,
			Transactions: []dto.Transaction{},
		},
		"2": {
			ID:           2,
			Balance:      0,
			Limit:        80000,
			Transactions: []dto.Transaction{},
		},
		"3": {
			ID:           3,
			Balance:      0,
			Limit:        1000000,
			Transactions: []dto.Transaction{},
		},
		"4": {
			ID:           4,
			Balance:      0,
			Limit:        100000000,
			Transactions: []dto.Transaction{},
		},
		"5": {
			ID:           5,
			Balance:      0,
			Limit:        500000,
			Transactions: []dto.Transaction{},
		},
	}
}

func (c *Database) set(table map[string]dto.Account) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//c.Table[name]++
	c.Table = table
}

func main() {
	defer dbPql.Close()

	fromAsyncQueue := new(FromAsyncQueue)
	iSyncQ = queue.NewEvent(fromAsyncQueue)
	go iSyncQ.Consume()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.SetResponseJsonHeaders)
	r.Use(middleware.TimeoutLimit)
	r.Get("/clientes/{id}/extrato", GetTransactions)
	r.Post("/clientes/{id}/transacoes", CreateTransaction)

	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	account := dbCache.Table[userID]
	if account.ID == 0 {
		http.Error(w, "Not found account", http.StatusNotFound)
		return
	}

	output := dto.GetBalanceOutputDTO{
		Ballance: dto.Ballance{
			Total: int(account.Balance),
			Date:  time.Now(),
			Limit: int(account.Limit),
		},
		Transactions: account.Transactions,
	}

	w.WriteHeader(http.StatusOK)
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

	account := dbCache.Table[accountIDParam]
	if account.ID == 0 {
		http.Error(w, "Not found account", http.StatusNotFound)
		return
	}

	if input.Type == "c" {
		account.Balance += input.Value
	}

	if input.Type == "d" {
		if account.Balance+account.Limit >= input.Value {
			http.Error(w, "Debit must be positive", http.StatusUnprocessableEntity)
			return
		}
		account.Balance -= input.Value
	}

	input.CreatedAt = time.Now()
	if len(account.Transactions) == 5 {
		account.Transactions = account.Transactions[1:5]
		account.Transactions = append(account.Transactions, input)
	} else {
		account.Transactions = append(account.Transactions, input)
	}

	b, _ := json.Marshal(account)
	iSyncQ.Publish(b)

	output = dto.TransactionOutput{
		Limit:   account.Limit,
		Balance: account.Balance,
	}

	dbCache.set(map[string]dto.Account{
		accountIDParam: account,
	})

	//w.WriteHeader(http.StatusOK)//todo: suspeita de estar dando superflus headers
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

type FromAsyncQueue struct{}

func (*FromAsyncQueue) Consumer(input []byte) error {
	var account dto.Account
	if err := json.Unmarshal(input, &account); err != nil {
		return err
	}

	const query = `INSERT INTO "transactions" ("id", "value", "type", "description", "client_id") VALUES ($1, $2, $3, $4, $5);`
	_, err := dbPql.Exec(
		query,
		uuid.NewString(),
		account.Transactions[0].Value,
		account.Transactions[0].Type,
		account.Transactions[0].Desc,
		account.ID,
	)

	return err
}

func (q FromAsyncQueue) ConsumerErr(input queue.Retry) error {
	logger.Error("The Consumer Error", "error", input.RetrieveError, "data", string(input.Msg), "qtd", input.Quantity)

	var account dto.Account
	if err := json.Unmarshal(input.Msg, &account); err != nil {
		return err
	}

	const query = `INSERT INTO "transactions" ("id", "value", "type", "description", "client_id") VALUES ($1, $2, $3, $4, $5);`
	_, err := dbPql.Exec(
		query,
		uuid.NewString(),
		account.Transactions[0].Value,
		account.Transactions[0].Type,
		account.Transactions[0].Desc,
		account.ID,
	)
	return err
}
