package delete_client

import (
	"database/sql"
	"errors"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type DeleteClientOutputDTO struct {
	ID string `json:"id"`
}

type DeleteClientUseCase struct {
	ClientGateway  gateway.ClientGateway
	AccountGateway gateway.AccountGateway
}

func NewDeleteClientUseCase(clientGateway gateway.ClientGateway, accountGateway gateway.AccountGateway) *DeleteClientUseCase {
	return &DeleteClientUseCase{
		ClientGateway:  clientGateway,
		AccountGateway: accountGateway,
	}
}

func (u *DeleteClientUseCase) Execute(id string) (*DeleteClientOutputDTO, error) {
	if id == "" {
		return nil, apperror.NewValidation("client id is required")
	}

	_, err := u.ClientGateway.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("client not found")
		}

		return nil, apperror.NewInternal("failed to get client", err)
	}

	_, err = u.AccountGateway.FindByClientID(id)
	if err == nil {
		return nil, apperror.NewConflict("client has linked accounts")
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NewInternal("failed to check client accounts", err)
	}

	err = u.ClientGateway.Delete(id)
	if err != nil {
		return nil, apperror.NewInternal("failed to delete client", err)
	}

	return &DeleteClientOutputDTO{ID: id}, nil
}
