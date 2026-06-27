package createtransaction

import (
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/events"
)

type CreateTransactionInputDTO struct {
	AccountFromID string
	AccountToID   string
	Amount        float64
}

type CreateTransactionOutputDTO struct {
	ID string
}

type CreateTransactionUseCase struct {
	TransactionGateway gateway.TransactionGateway
	AccountGateway     gateway.AccountGateway
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
}

func NewCreateTransactionUseCase(
	transactionGateway gateway.TransactionGateway,
	accountGateway gateway.AccountGateway,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreatedEvent events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		TransactionGateway: transactionGateway,
		AccountGateway:     accountGateway,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreatedEvent,
	}
}

func (u *CreateTransactionUseCase) Execute(input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	accountFrom, err := u.AccountGateway.FindByID(input.AccountFromID)
	if err != nil {
		return nil, err
	}

	accountTo, err := u.AccountGateway.FindByID(input.AccountToID)
	if err != nil {
		return nil, err
	}

	transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
	if err != nil {
		return nil, err
	}

	err = u.TransactionGateway.Create(transaction)
	if err != nil {
		return nil, err
	}

	u.TransactionCreated.SetPayload(transaction)
	u.EventDispatcher.Dispatch(u.TransactionCreated)

	return &CreateTransactionOutputDTO{
		ID: transaction.ID,
	}, nil
}
