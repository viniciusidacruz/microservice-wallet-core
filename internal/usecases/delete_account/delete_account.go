package delete_account

import (
	"database/sql"
	"errors"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type DeleteAccountOutputDTO struct {
	ID string `json:"id"`
}

type DeleteAccountUseCase struct {
	AccountGateway     gateway.AccountGateway
	TransactionGateway gateway.TransactionGateway
}

func NewDeleteAccountUseCase(accountGateway gateway.AccountGateway, transactionGateway gateway.TransactionGateway) *DeleteAccountUseCase {
	return &DeleteAccountUseCase{
		AccountGateway:     accountGateway,
		TransactionGateway: transactionGateway,
	}
}

func (u *DeleteAccountUseCase) Execute(id string) (*DeleteAccountOutputDTO, error) {
	if id == "" {
		return nil, apperror.NewValidation("account id is required")
	}

	_, err := u.AccountGateway.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("account not found")
		}

		return nil, apperror.NewInternal("failed to get account", err)
	}

	hasTransactions, err := u.TransactionGateway.ExistsByAccountID(id)
	if err != nil {
		return nil, apperror.NewInternal("failed to check account transactions", err)
	}

	if hasTransactions {
		return nil, apperror.NewConflict("account has linked transactions")
	}

	err = u.AccountGateway.Delete(id)
	if err != nil {
		return nil, apperror.NewInternal("failed to delete account", err)
	}

	return &DeleteAccountOutputDTO{ID: id}, nil
}
