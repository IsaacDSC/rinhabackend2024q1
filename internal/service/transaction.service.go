package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"rinhabackend/internal/entity"
	"rinhabackend/internal/repository"
	"rinhabackend/shared/dto"
	"rinhabackend/shared/interfaces"
	"time"
)

type TransactionService struct {
	transactionRepository interfaces.TransactionRepository
	clientRepository      *repository.ClientRepository //TODO: mudar para interface
	db                    *sql.DB
}

func NewTransactionService(
	transactionRepository interfaces.TransactionRepository,
	clientRepository *repository.ClientRepository, //TODO: mudar para interface
	db *sql.DB,
) *TransactionService {
	return &TransactionService{
		transactionRepository,
		clientRepository,
		db,
	}
}

func (ts *TransactionService) CreateTransaction(ctx context.Context, input *dto.TransactionInput, userID string) (output dto.TransactionOutput, err error) {

	tx, err := ts.db.Begin()
	if err != nil {
		return
	}

	clients, err := ts.clientRepository.GetClients(ctx, tx)
	if err != nil {
		return
	}

	var user struct {
		id             string
		limit, balance int64
	}

	for i := range clients {
		if clients[i].ID == userID {
			user.id = clients[i].ID
			user.limit = clients[i].Limit
			user.balance = clients[i].Balance
			break
		}
	}

	if user.id == "" {
		err = errors.New(string(entity.ERROR_CLIENT_NOT_FOUND))
		return
	}

	transaction := input.ToDomain(userID, user.limit, user.balance)
	if err = transaction.Execute(); err != nil {
		return
	}

	if err = ts.transactionRepository.CreateTransaction(
		ctx,
		*transaction,
	); err != nil {
		return
	}

	fmt.Println("transaction.Balance", transaction.Balance)
	err = ts.clientRepository.UpdateBalancer(ctx, tx, transaction.UserID, transaction.Balance)

	return
}

func (ts *TransactionService) GetTransaction(ctx context.Context, userID string) (output []dto.TransactionsOutput, err error) {
	transactions, err := ts.transactionRepository.GetTransactionsByUser(ctx, userID)
	if err != nil {
		return
	}

	if len(transactions) == 0 {
		return
	}

	balance := dto.Ballance{
		Total: int(transactions[0].Balance),
		Date:  time.Now(),
		Limit: int(transactions[0].Limit),
	}

	var tx []dto.Transactions
	for i := range transactions {
		tx = append(tx, dto.Transactions{
			Value:     int(transactions[i].Value),
			Type:      string(transactions[i].Type),
			Desc:      transactions[i].Description,
			CreatedAt: transactions[i].CreatedAt,
		})
	}

	output = append(output, dto.TransactionsOutput{
		Ballance:     balance,
		Transactions: tx,
	})

	return
}
