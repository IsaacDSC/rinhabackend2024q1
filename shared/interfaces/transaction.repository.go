package interfaces

import (
	"context"
	"rinhabackend/internal/entity"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction entity.Transaction) error
	GetTransactionsByUser(ctx context.Context, userID string) ([]entity.Transaction, error)
}
