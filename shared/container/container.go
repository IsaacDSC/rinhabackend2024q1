package container

import (
	"database/sql"
	"rinhabackend/internal/repository"
	"rinhabackend/internal/service"
	"rinhabackend/web/handler"
)

type ContainerRepositories struct {
	ClientRepository      *repository.ClientRepository
	TransactionRepository *repository.TransactionRepository
}

func NewContainerRepositories(db *sql.DB) *ContainerRepositories {
	return &ContainerRepositories{
		ClientRepository:      repository.NewClientRepository(db),
		TransactionRepository: repository.NewTransactionRepository(db),
	}
}

type ContainerServices struct {
	TransactionService *service.TransactionService //TODO: mudar para interface
}

func NewContainerService(repo *ContainerRepositories, db *sql.DB) *ContainerServices {
	return &ContainerServices{
		TransactionService: service.NewTransactionService(repo.TransactionRepository, repo.ClientRepository, db),
	}
}

type ContainerHandlers struct {
	TransactionHandler *handler.TransactionHandler
}

func NewContainerHandlers(containerServices *ContainerServices) *ContainerHandlers {
	return &ContainerHandlers{
		TransactionHandler: handler.NewTransactionHandler(containerServices.TransactionService),
	}
}
