package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	CREDIT_TRANSACTION_TYPE TransactionType = "c"
	DEBIT_TRANSACTION_TYPE  TransactionType = "d"
)

type Transaction struct {
	ID          uuid.UUID
	UserID      string
	Value       int64
	Type        TransactionType
	Description string
	Limit       int64
	Balance     int64
	CreatedAt   time.Time
}

func (t *Transaction) Execute() (err error) {
	if t.Type == DEBIT_TRANSACTION_TYPE {
		if err := t.debit(); err != nil {
			return err
		}
		return
	}

	if t.Type == CREDIT_TRANSACTION_TYPE {
		if err := t.credit(); err != nil {
			return err
		}
		return
	}

	return errors.New(string(ERROR_INVALID_OPERATION))
}

func (t *Transaction) credit() (err error) {
	t.Balance = t.Balance + t.Value
	return
}

func (t *Transaction) debit() (err error) {
	if t.Balance <= 0 && t.Value <= t.Limit {
		return errors.New("credits must be positive")
	}
	t.Balance = t.Value - t.Balance
	return
}
