package get_client

import (
	"database/sql"
	"errors"
	"time"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type GetClientOutputDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetClientUseCase struct {
	ClientGateway gateway.ClientGateway
}

func NewGetClientUseCase(clientGateway gateway.ClientGateway) *GetClientUseCase {
	return &GetClientUseCase{
		ClientGateway: clientGateway,
	}
}

func (u *GetClientUseCase) Execute(id string) (*GetClientOutputDTO, error) {
	if id == "" {
		return nil, apperror.NewValidation("client id is required")
	}

	client, err := u.ClientGateway.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("client not found")
		}

		return nil, apperror.NewInternal("failed to get client", err)
	}

	return &GetClientOutputDTO{
		ID:        client.ID,
		Name:      client.Name,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}, nil
}
