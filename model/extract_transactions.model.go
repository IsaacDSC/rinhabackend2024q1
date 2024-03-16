package model

import (
	"database/sql"
)

type ExtractTransactionModel struct {
	TransactionID    sql.NullString
	ClientID         sql.NullString
	TransactionType  sql.NullString
	TransactionDesc  sql.NullString
	TransactionValue sql.NullInt64
	CreatedAt        sql.NullTime
	AccountBalance   sql.NullInt64
	AccountLimit     sql.NullInt64
}
