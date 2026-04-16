package gateway

import "github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"

type AccountGateway interface {
	FindByID(id string) (*entity.Account, error)
	Save(*entity.Account) error
}