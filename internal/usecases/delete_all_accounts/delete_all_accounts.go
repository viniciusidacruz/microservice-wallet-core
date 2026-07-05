package delete_all_accounts

import (
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type DeleteAllAccountsOutputDTO struct {
	DeletedCount int64 `json:"deleted_count"`
}

type DeleteAllAccountsUseCase struct {
	AccountGateway     gateway.AccountGateway
	TransactionGateway gateway.TransactionGateway
}

func NewDeleteAllAccountsUseCase(accountGateway gateway.AccountGateway, transactionGateway gateway.TransactionGateway) *DeleteAllAccountsUseCase {
	return &DeleteAllAccountsUseCase{
		AccountGateway:     accountGateway,
		TransactionGateway: transactionGateway,
	}
}

func (u *DeleteAllAccountsUseCase) Execute() (*DeleteAllAccountsOutputDTO, error) {
	if _, err := u.TransactionGateway.DeleteAll(); err != nil {
		return nil, apperror.NewInternal("failed to delete transactions", err)
	}

	deletedCount, err := u.AccountGateway.DeleteAll()
	if err != nil {
		return nil, apperror.NewInternal("failed to delete accounts", err)
	}

	return &DeleteAllAccountsOutputDTO{DeletedCount: deletedCount}, nil
}
