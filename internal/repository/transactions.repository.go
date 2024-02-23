package repository

import (
	"context"
	"database/sql"
	"rinhabackend/internal/entity"
	"rinhabackend/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db}
}

func (t *TransactionRepository) CreateTransaction(
	ctx context.Context, tx *sql.Tx, transaction entity.Transaction,
) error {
	const query = `
	INSERT INTO "transactions" ("id", "value", "type", "description", "client_id") VALUES ($1, $2, $3, $4, $5);`
	if _, err := tx.Exec(
		query,
		transaction.ID,
		transaction.Value,
		transaction.Type,
		transaction.Description,
		transaction.Client.ID,
	); err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepository) GetTransactionsByUser(ctx context.Context, clientID string) (output []entity.Transaction, err error) {
	const query = `
	select 
	clients.balance,
	clients.limit,
	clients.id,
	transactions.id,
	transactions.value,
	transactions.type,
	transactions.description,
	transactions.created_at 
	from clients 
	left join transactions on clients.id = "transactions".client_id 
	where clients.id = $1  ORDER BY "created_at" DESC limit 10;
	`

	rows, err := t.db.Query(query, clientID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t models.TransactionModel
		if err = rows.Scan(
			&t.Client.Balance,
			&t.Client.Limit,
			&t.ClientID,
			&t.ID,
			&t.Value,
			&t.Type,
			&t.Description,
			&t.CreatedAt,
		); err != nil {
			return
		}
		output = append(output, entity.Transaction{
			ID:          t.ID.String,
			Value:       t.Value.Int64,
			Type:        entity.TransactionType(t.Type.String),
			Description: t.Description.String,
			CreatedAt:   t.CreatedAt.Time,
			Client: entity.Client{
				ID:      int(t.ClientID.Int16),
				Balance: t.Client.Balance.Int64,
				Limit:   t.Client.Limit.Int64,
			},
		})
	}

	return
}
