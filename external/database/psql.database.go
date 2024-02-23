package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func ConnPSQL() *sql.DB {
	dbPort := os.Getenv("DB_PORT")
	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://root:root@192.168.1.100:%s/rinha_backend?sslmode=disable", dbPort))
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(83)
	db.SetMaxIdleConns(20)
	return db
}
