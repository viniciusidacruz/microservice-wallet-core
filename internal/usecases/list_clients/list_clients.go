package list_clients

import (
	"time"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type ClientItemDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListClientsOutputDTO struct {
	Clients []ClientItemDTO `json:"clients"`
}

type ListClientsUseCase struct {
	ClientGateway gateway.ClientGateway
}

func NewListClientsUseCase(clientGateway gateway.ClientGateway) *ListClientsUseCase {
	return &ListClientsUseCase{
		ClientGateway: clientGateway,
	}
}

func (u *ListClientsUseCase) Execute() (*ListClientsOutputDTO, error) {
	clients, err := u.ClientGateway.List()
	if err != nil {
		return nil, apperror.NewInternal("failed to list clients", err)
	}

	items := make([]ClientItemDTO, 0, len(clients))
	for _, client := range clients {
		items = append(items, ClientItemDTO{
			ID:        client.ID,
			Name:      client.Name,
			Email:     client.Email,
			CreatedAt: client.CreatedAt,
			UpdatedAt: client.UpdatedAt,
		})
	}

	return &ListClientsOutputDTO{
		Clients: items,
	}, nil
}
