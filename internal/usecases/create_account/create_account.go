package create_account

import (
	"database/sql"
	"errors"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type CreateAccountInputDTO struct {
	ClientID string `json:"client_id"`
}

type CreateAccountOutputDTO struct {
	ID string `json:"id"`
}

type CreateAccountUseCase struct {
	AccountGateway gateway.AccountGateway
	ClientGateway  gateway.ClientGateway
}

func NewCreateAccountUseCase(accountGateway gateway.AccountGateway, clientGateway gateway.ClientGateway) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		AccountGateway: accountGateway,
		ClientGateway:  clientGateway,
	}
}

func (u *CreateAccountUseCase) Execute(input CreateAccountInputDTO) (*CreateAccountOutputDTO, error) {
	if input.ClientID == "" {
		return nil, apperror.NewValidation("client_id is required")
	}

	client, err := u.ClientGateway.Get(input.ClientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("client not found")
		}

		return nil, apperror.NewInternal("failed to get client", err)
	}

	_, err = u.AccountGateway.FindByClientID(client.ID)
	if err == nil {
		return nil, apperror.NewConflict("client already has an account")
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NewInternal("failed to check client account", err)
	}

	account := entity.NewAccount(client)
	if account == nil {
		return nil, apperror.NewInternal("failed to create account", nil)
	}

	err = u.AccountGateway.Save(account)
	if err != nil {
		return nil, apperror.NewInternal("failed to save account", err)
	}

	return &CreateAccountOutputDTO{
		ID: account.ID,
	}, nil
}
