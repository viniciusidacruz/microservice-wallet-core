package uow

import (
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type PassThroughUnitOfWork struct {
	account     gateway.AccountGateway
	transaction gateway.TransactionGateway
}

func NewPassThroughUnitOfWork(accountGateway gateway.AccountGateway, transactionGateway gateway.TransactionGateway) *PassThroughUnitOfWork {
	return &PassThroughUnitOfWork{
		account:     accountGateway,
		transaction: transactionGateway,
	}
}

func (p *PassThroughUnitOfWork) Account() gateway.AccountGateway {
	return p.account
}

func (p *PassThroughUnitOfWork) Transaction() gateway.TransactionGateway {
	return p.transaction
}

func (p *PassThroughUnitOfWork) Do(fn func(repos Repositories) error) error {
	return fn(Repositories{
		Account:     p.account,
		Transaction: p.transaction,
	})
}
