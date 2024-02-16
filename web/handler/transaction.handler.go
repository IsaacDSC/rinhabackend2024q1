package handler

import (
	"encoding/json"
	"net/http"
	"rinhabackend/internal/service"
	"rinhabackend/shared/dto"

	"github.com/go-chi/chi"
)

type TransactionHandler struct {
	service *service.TransactionService //TODO: mudar para interface
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

func (th *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	output, err := th.service.GetTransaction(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (th *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	defer r.Body.Close()

	var (
		input  dto.TransactionInput
		output dto.TransactionOutput
		err    error
	)

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err = th.service.CreateTransaction(r.Context(), &input, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// var balance sql.NullInt64
	// if err := th.db.QueryRow(`SELECT balance FROM "transactions" WHERE id = $1;`, userID).Scan(
	// 	&balance,
	// ); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// if !balance.Valid {
	// 	http.Error(w, "balancer invalid", http.StatusInternalServerError)
	// 	return
	// }

	// if input.Type == entity.DEBIT_TRANSACTION_TYPE {
	// 	if balance.Int64 < input.Value {
	// 		http.Error(w, "transaction invalid", http.StatusBadRequest)
	// 		return
	// 	}
	// 	balance.Int64 = balance.Int64 - input.Value
	// }

	// if err := th.transactionRepository.CreateTransaction(userID, input.Value, input.Type, input.Description); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
