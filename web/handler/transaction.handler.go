package handler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"rinhabackend/internal/entity"
	"rinhabackend/internal/service"
	"rinhabackend/shared/dto"
	"strconv"
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
		return
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

	if input.Value <= 0 {
		http.Error(w, "Invalid value", http.StatusUnprocessableEntity)
		return
	}
	if input.Type != "c" && input.Type != "d" {
		http.Error(w, "Invalid type", http.StatusUnprocessableEntity)
		return
	}
	if len(input.Description) == 0 || len(input.Description) > 10 {
		http.Error(w, "Invalid Description", http.StatusUnprocessableEntity)
		return
	}

	clientID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err = th.service.CreateTransaction(r.Context(), &input, clientID)
	if err != nil {
		switch err.Error() {
		case string(entity.ERROR_CLIENT_NOT_FOUND):
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case string(entity.ERROR_DEBIT_MUST_BE_POSITIVE):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
