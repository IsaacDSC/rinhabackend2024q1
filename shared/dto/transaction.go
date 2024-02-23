package dto

import (
	"rinhabackend/internal/entity"
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

func (t *TransactionInput) ToDomain(transactionID string, clientID int, limit, balance int64) *entity.Transaction {
	return &entity.Transaction{
		ID:          transactionID,
		Value:       t.Value,
		Type:        t.Type,
		Description: t.Description,
		Client: entity.Client{
			ID:      clientID,
			Balance: balance,
			Limit:   limit,
		},
	}
}
