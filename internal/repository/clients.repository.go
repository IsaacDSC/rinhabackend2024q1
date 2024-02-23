package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"rinhabackend/shared/dto"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db}
}

func (c *ClientRepository) GetClient(ctx context.Context, tx *sql.Tx, clientID int) (output dto.Client, err error) {
	err = tx.QueryRow(`SELECT "id", "limit", "balance" FROM "clients" where id = $1 FOR UPDATE;`, clientID).Scan(
		&output.ID,
		&output.Limit,
		&output.Balance,
	)
	return
}

func (c *ClientRepository) UpdateBalancer(ctx context.Context, tx *sql.Tx, userID int, balancer int64) error {
	if _, err := tx.Exec(`UPDATE "clients" SET "balance" = $1 WHERE "id" = $2`, balancer, userID); err != nil {
		return err
	}
	return nil
}

func (c *ClientRepository) CreateCacheClients(clients []dto.Client) {
	b, err := json.Marshal(clients)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("./tmp/cache/clients.json", b, 0644); err != nil {
		log.Fatal(err)
	}
}
