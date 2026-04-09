package gateway

import "github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(*entity.Client) error
}