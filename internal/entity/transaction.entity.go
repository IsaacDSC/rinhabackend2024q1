package entity

import (
	"errors"
	"time"
)

type TransactionType string

const (
	CREDIT_TRANSACTION_TYPE TransactionType = "c"
	DEBIT_TRANSACTION_TYPE  TransactionType = "d"
)

type Transaction struct {
	ID          string
	Value       int64
	Type        TransactionType
	Description string
	CreatedAt   time.Time
	Client      Client
}

type Client struct {
	ID      int
	Balance int64
	Limit   int64
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
	t.Client.Balance += t.Value
	return
}

func (t *Transaction) debit() (err error) {
	// if t.Value > t.Balance || t.Value > t.Limit {
	// 	return errors.New("debit must be positive")
	// }
	if t.Client.Balance+t.Client.Limit >= t.Value {
		return errors.New(string(ERROR_DEBIT_MUST_BE_POSITIVE))
	}
	t.Client.Balance -= t.Value
	return
}
