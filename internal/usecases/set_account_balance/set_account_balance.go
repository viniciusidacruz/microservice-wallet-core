package set_account_balance

import (
	"database/sql"
	"errors"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type SetAccountBalanceInputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type SetAccountBalanceOutputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type SetAccountBalanceUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewSetAccountBalanceUseCase(accountGateway gateway.AccountGateway) *SetAccountBalanceUseCase {
	return &SetAccountBalanceUseCase{
		AccountGateway: accountGateway,
	}
}

func (u *SetAccountBalanceUseCase) Execute(input SetAccountBalanceInputDTO) (*SetAccountBalanceOutputDTO, error) {
	if input.AccountID == "" {
		return nil, apperror.NewValidation("account_id is required")
	}

	if input.Balance < 0 {
		return nil, apperror.NewValidation("balance must be greater than or equal to zero")
	}

	account, err := u.AccountGateway.FindByID(input.AccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("account not found")
		}

		return nil, apperror.NewInternal("failed to get account", err)
	}

	err = account.SetBalance(input.Balance)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	err = u.AccountGateway.UpdateBalance(account)
	if err != nil {
		return nil, apperror.NewInternal("failed to update account balance", err)
	}

	return &SetAccountBalanceOutputDTO{
		AccountID: account.ID,
		Balance:   account.Balance,
	}, nil
}
