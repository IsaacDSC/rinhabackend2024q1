package repository

import (
	"context"
	"database/sql"
	"rinhabackend/internal/entity"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db}
}

func (t *TransactionRepository) CreateTransaction(
	ctx context.Context, transaction entity.Transaction,
) error {
	const query = `
	INSERT INTO "transactions" ("id","user_id", "value", "type", "description", "limit", "balance") VALUES ($1, $2, $3, $4, $5, $6, $7);`
	if _, err := t.db.Exec(
		query,
		transaction.ID,
		transaction.UserID,
		transaction.Balance,
		transaction.Type,
		transaction.Description,
		transaction.Limit,
		transaction.Balance,
	); err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepository) GetTransactionsByUser(ctx context.Context, userID string) (output []entity.Transaction, err error) {
	const query = `SELECT "id", "user_id", "value", "type", "description", "balance", "limit", "created_at" FROM "transactions" WHERE user_id = $1 ORDER BY "created_at" DESC;`
	rows, err := t.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t entity.Transaction
		if err = rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Value,
			&t.Type,
			&t.Description,
			&t.Balance,
			&t.Limit,
			&t.CreatedAt,
		); err != nil {
			return
		}
		output = append(output, t)
	}

	return
}
