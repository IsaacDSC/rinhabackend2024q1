package dto

import (
	"rinhabackend/internal/entity"

	"github.com/google/uuid"
)

type TransactionInput struct {
	Value       int64                  `json:"valor"`
	Type        entity.TransactionType `json:"tipo"`
	Description string                 `json:"descricao"`
}

type TransactionOutput struct {
	Limit   int64 `json:"limite"`
	Balance int64 `json:"saldo"`
}

func (t *TransactionInput) ToDomain(userID string, limit, balance int64) *entity.Transaction {
	return &entity.Transaction{
		ID:          uuid.New(),
		UserID:      userID,
		Value:       t.Value,
		Type:        t.Type,
		Description: t.Description,
		Balance:     balance,
		Limit:       limit,
	}
}
