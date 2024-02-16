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

func (c *ClientRepository) GetClients(ctx context.Context, tx *sql.Tx) (output []dto.Client, err error) {
	rows, err := c.db.Query(`SELECT "id", "limit", "balance" FROM "clients";`)
	if err != nil {
		return
	}

	for rows.Next() {
		var c dto.Client
		if err = rows.Scan(
			&c.ID,
			&c.Limit,
			&c.Balance,
		); err != nil {
			return
		}
		output = append(output, c)
	}

	return
}

func (c *ClientRepository) UpdateBalancer(ctx context.Context, tx *sql.Tx, userID string, balancer int64) (err error) {
	_, err = c.db.Exec(`UPDATE "transactions" SET "balance" = $1 WHERE "id" = $2`, balancer, userID)
	return
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
