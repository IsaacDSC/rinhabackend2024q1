package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func ConnPSQL() *sql.DB {
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	connectionURL := fmt.Sprintf("postgresql://root:root@%s:%s/rinha_backend?sslmode=disable", dbHost, dbPort)
	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(83)
	db.SetMaxIdleConns(20)
	return db
}
