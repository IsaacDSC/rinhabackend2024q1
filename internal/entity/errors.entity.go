package entity

type ERRORS string

const (
	ERROR_INVALID_OPERATION      ERRORS = "Invalid operation"
	ERROR_CLIENT_NOT_FOUND       ERRORS = "Client not found"
	ERROR_DEBIT_MUST_BE_POSITIVE ERRORS = "Debit must be positive"
)
