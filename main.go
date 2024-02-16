package main

import (
	"database/sql"
	"log"
	"net/http"
	"rinhabackend/external/database"
	"rinhabackend/shared/container"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

var (
	db *sql.DB
)

func init() {
	db = database.ConnPSQL()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	repositories := container.NewContainerRepositories(db)
	services := container.NewContainerService(repositories, db)
	handlers := container.NewContainerHandlers(services)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/clientes/{id}/extrato", handlers.TransactionHandler.GetTransactions)
	r.Post("/clientes/{id}/transacoes", handlers.TransactionHandler.CreateTransaction)

	http.ListenAndServe(":3000", r)
}
