package interfaces

import (
	"context"
	"database/sql"
	"rinhabackend/internal/entity"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *sql.Tx, transaction entity.Transaction) error
	GetTransactionsByUser(ctx context.Context, userID string) ([]entity.Transaction, error)
}
