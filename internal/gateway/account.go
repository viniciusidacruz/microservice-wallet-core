package gateway

import "github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"

type AccountGateway interface {
	FindByID(id string) (*entity.Account, error)
	FindByClientID(clientID string) (*entity.Account, error)
	List() ([]*entity.Account, error)
	Save(*entity.Account) error
	UpdateBalance(account *entity.Account) error
	Delete(id string) error
	DeleteAll() (int64, error)
}