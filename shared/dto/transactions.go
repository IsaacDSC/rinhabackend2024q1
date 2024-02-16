package dto

import "time"

type TransactionsOutput struct {
	Ballance     Ballance       `json:"saldo"`
	Transactions []Transactions `json:"ultimas_transacoes"`
}

type Ballance struct {
	Total int       `json:"total"`
	Date  time.Time `json:"data_extrato"`
	Limit int       `json:"limite"`
}

type Transactions struct {
	Value     int       `json:"valor"`
	Type      string    `json:"tipo"`
	Desc      string    `json:"descricao"`
	CreatedAt time.Time `json:"realizada_em"`
}
