package database

import (
	"database/sql"
	"log"
)

func ConnPSQL() *sql.DB {
	db, err := sql.Open("postgres", "postgresql://root:root@192.168.1.100:5432/rinha_backend?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	return db
}
