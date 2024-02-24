package dto

import "time"

type TransactionOutput struct {
	Limit   int64 `json:"limite"`
	Balance int64 `json:"saldo"`
}

type GetBalanceOutputDTO struct {
	Ballance     Ballance      `json:"saldo"`
	Transactions []Transaction `json:"ultimas_transacoes"`
}

type Ballance struct {
	Total int       `json:"total"`
	Limit int       `json:"limite"`
	Date  time.Time `json:"data_extrato"`
}

type Account struct {
	ID           int32
	Balance      int64
	Limit        int64
	Transactions []Transaction
}

type Transaction struct {
	Value     int64     `json:"valor"`
	Type      string    `json:"tipo"`
	Desc      string    `json:"descricao"`
	CreatedAt time.Time `json:"realizada_em"`
}
