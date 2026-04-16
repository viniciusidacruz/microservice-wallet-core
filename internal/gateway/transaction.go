package gateway

import "github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}