package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"rinhabackend/internal/entity"
	"rinhabackend/internal/repository"
	"rinhabackend/shared/dto"
	"rinhabackend/shared/interfaces"
	"time"

	"github.com/google/uuid"
)

type TransactionService struct {
	transactionRepository interfaces.TransactionRepository
	clientRepository      *repository.ClientRepository //TODO: mudar para interface
	db                    *sql.DB
	logger                *slog.Logger
}

func NewTransactionService(
	transactionRepository interfaces.TransactionRepository,
	clientRepository *repository.ClientRepository, //TODO: mudar para interface
	db *sql.DB,
	logger *slog.Logger,
) *TransactionService {
	return &TransactionService{
		transactionRepository,
		clientRepository,
		db,
		logger,
	}
}

func (ts *TransactionService) CreateTransaction(ctx context.Context, input *dto.TransactionInput, clientID int) (output dto.TransactionOutput, err error) {
	transactionID := uuid.NewString()
	tx, err := ts.db.Begin()
	if err != nil {
		ts.logger.Error(fmt.Errorf("Error creating transaction: %v", err).Error(), "transaction_id", transactionID)
		return
	}
	defer tx.Rollback()

	client, err := ts.clientRepository.GetClient(ctx, tx, clientID)
	if err != nil {
		ts.logger.Error(fmt.Errorf("Error get client: %v", err).Error(), "transaction_id", transactionID)
		return
	}

	if client.ID == 0 {
		err = errors.New(string(entity.ERROR_CLIENT_NOT_FOUND))
		return
	}

	transaction := input.ToDomain(transactionID, clientID, client.Limit, client.Balance)
	if err = transaction.Execute(); err != nil {
		return
	}

	if err = ts.transactionRepository.CreateTransaction(
		ctx,
		tx,
		*transaction,
	); err != nil {
		ts.logger.Error(fmt.Errorf("Create transaction error save to db: %v", err).Error(), "transaction_id", transactionID)
		return
	}

	if err = ts.clientRepository.UpdateBalancer(ctx, tx, int(transaction.Client.ID), transaction.Client.Balance); err != nil {
		ts.logger.Error(fmt.Errorf("Update balancer with error save to db: %v", err).Error(), "transaction_id", transactionID)
		return
	}

	output = dto.TransactionOutput{
		Limit:   transaction.Client.Limit,
		Balance: transaction.Client.Balance,
	}

	err = tx.Commit()
	return
}

func (ts *TransactionService) GetTransaction(ctx context.Context, userID string) (output dto.TransactionsOutput, err error) {
	transactions, err := ts.transactionRepository.GetTransactionsByUser(ctx, userID)
	if err != nil {
		ts.logger.Error(fmt.Errorf("Get ticket the transactions with error: %v", err).Error(), "client_id", userID)
		return
	}

	if len(transactions) == 0 {
		return
	}

	balance := dto.Ballance{
		Total: int(transactions[0].Client.Balance),
		Date:  time.Now(),
		Limit: int(transactions[0].Client.Limit),
	}

	var tx []dto.Transactions
	for i := range transactions {
		if transactions[i].ID == "" {
			tx = []dto.Transactions{}
			break
		}
		tx = append(tx, dto.Transactions{
			Value:     int(transactions[i].Value),
			Type:      string(transactions[i].Type),
			Desc:      transactions[i].Description,
			CreatedAt: transactions[i].CreatedAt,
		})
	}

	output = dto.TransactionsOutput{
		Ballance:     balance,
		Transactions: tx,
	}

	return
}
