package delete_all_clients

import (
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type DeleteAllClientsOutputDTO struct {
	DeletedCount int64 `json:"deleted_count"`
}

type DeleteAllClientsUseCase struct {
	ClientGateway      gateway.ClientGateway
	AccountGateway     gateway.AccountGateway
	TransactionGateway gateway.TransactionGateway
}

func NewDeleteAllClientsUseCase(
	clientGateway gateway.ClientGateway,
	accountGateway gateway.AccountGateway,
	transactionGateway gateway.TransactionGateway,
) *DeleteAllClientsUseCase {
	return &DeleteAllClientsUseCase{
		ClientGateway:      clientGateway,
		AccountGateway:     accountGateway,
		TransactionGateway: transactionGateway,
	}
}

func (u *DeleteAllClientsUseCase) Execute() (*DeleteAllClientsOutputDTO, error) {
	if _, err := u.TransactionGateway.DeleteAll(); err != nil {
		return nil, apperror.NewInternal("failed to delete transactions", err)
	}

	if _, err := u.AccountGateway.DeleteAll(); err != nil {
		return nil, apperror.NewInternal("failed to delete accounts", err)
	}

	deletedCount, err := u.ClientGateway.DeleteAll()
	if err != nil {
		return nil, apperror.NewInternal("failed to delete clients", err)
	}

	return &DeleteAllClientsOutputDTO{DeletedCount: deletedCount}, nil
}
