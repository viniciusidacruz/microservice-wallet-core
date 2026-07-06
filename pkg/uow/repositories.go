package uow

import (
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type Repositories struct {
	Account     gateway.AccountGateway
	Transaction gateway.TransactionGateway
}

type UnitOfWork interface {
	Account() gateway.AccountGateway
	Transaction() gateway.TransactionGateway
	Do(fn func(repos Repositories) error) error
}
