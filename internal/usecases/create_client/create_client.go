package create_client

import (
	"database/sql"
	"errors"
	"time"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type CreateClientInputDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateClientOutputDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateClientUseCase struct {
	ClientGateway gateway.ClientGateway
}

func NewCreateClientUseCase(clientGateway gateway.ClientGateway) *CreateClientUseCase {
	return &CreateClientUseCase{
		ClientGateway: clientGateway,
	}
}

func (u *CreateClientUseCase) Execute(input CreateClientInputDTO) (*CreateClientOutputDTO, error) {
	client, err := entity.NewClient(input.Name, input.Email)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	_, err = u.ClientGateway.GetByEmail(client.Email)
	if err == nil {
		return nil, apperror.NewConflict("client with this email already exists")
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NewInternal("failed to check client email", err)
	}

	err = u.ClientGateway.Save(client)
	if err != nil {
		return nil, apperror.NewInternal("failed to save client", err)
	}

	return &CreateClientOutputDTO{
		ID:        client.ID,
		Name:      client.Name,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}, nil
}
