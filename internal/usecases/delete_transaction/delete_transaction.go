package delete_transaction

import (
	"database/sql"
	"errors"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type DeleteTransactionOutputDTO struct {
	ID string `json:"id"`
}

type DeleteTransactionUseCase struct {
	TransactionGateway gateway.TransactionGateway
}

func NewDeleteTransactionUseCase(transactionGateway gateway.TransactionGateway) *DeleteTransactionUseCase {
	return &DeleteTransactionUseCase{
		TransactionGateway: transactionGateway,
	}
}

func (u *DeleteTransactionUseCase) Execute(id string) (*DeleteTransactionOutputDTO, error) {
	if id == "" {
		return nil, apperror.NewValidation("transaction id is required")
	}

	_, err := u.TransactionGateway.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("transaction not found")
		}

		return nil, apperror.NewInternal("failed to get transaction", err)
	}

	err = u.TransactionGateway.Delete(id)
	if err != nil {
		return nil, apperror.NewInternal("failed to delete transaction", err)
	}

	return &DeleteTransactionOutputDTO{ID: id}, nil
}
