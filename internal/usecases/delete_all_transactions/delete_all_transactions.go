package delete_all_transactions

import (
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type DeleteAllTransactionsOutputDTO struct {
	DeletedCount int64 `json:"deleted_count"`
}

type DeleteAllTransactionsUseCase struct {
	TransactionGateway gateway.TransactionGateway
}

func NewDeleteAllTransactionsUseCase(transactionGateway gateway.TransactionGateway) *DeleteAllTransactionsUseCase {
	return &DeleteAllTransactionsUseCase{
		TransactionGateway: transactionGateway,
	}
}

func (u *DeleteAllTransactionsUseCase) Execute() (*DeleteAllTransactionsOutputDTO, error) {
	deletedCount, err := u.TransactionGateway.DeleteAll()
	if err != nil {
		return nil, apperror.NewInternal("failed to delete transactions", err)
	}

	return &DeleteAllTransactionsOutputDTO{DeletedCount: deletedCount}, nil
}
