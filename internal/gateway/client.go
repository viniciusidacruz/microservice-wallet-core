package gateway

import "github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	GetByEmail(email string) (*entity.Client, error)
	List() ([]*entity.Client, error)
	Save(*entity.Client) error
	Delete(id string) error
	DeleteAll() (int64, error)
}