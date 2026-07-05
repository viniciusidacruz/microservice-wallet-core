package create_transaction

import (
	"database/sql"
	"errors"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/events"
)

type CreateTransactionInputDTO struct {
	AccountFromID string  `json:"account_from_id"`
	AccountToID   string  `json:"account_to_id"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID string `json:"id"`
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
	if input.AccountFromID == "" {
		return nil, apperror.NewValidation("account_from_id is required")
	}

	if input.AccountToID == "" {
		return nil, apperror.NewValidation("account_to_id is required")
	}

	if input.AccountFromID == input.AccountToID {
		return nil, apperror.NewValidation("account_from_id and account_to_id must be different")
	}

	if input.Amount <= 0 {
		return nil, apperror.NewValidation("amount must be greater than zero")
	}

	accountFrom, err := u.AccountGateway.FindByID(input.AccountFromID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("account from not found")
		}

		return nil, apperror.NewInternal("failed to get account from", err)
	}

	accountTo, err := u.AccountGateway.FindByID(input.AccountToID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("account to not found")
		}

		return nil, apperror.NewInternal("failed to get account to", err)
	}

	transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
	if err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	exists, err := u.TransactionGateway.Exists(input.AccountFromID, input.AccountToID, input.Amount)
	if err != nil {
		return nil, apperror.NewInternal("failed to check transaction duplicate", err)
	}

	if exists {
		return nil, apperror.NewConflict("transaction already exists")
	}

	err = u.TransactionGateway.Create(transaction)
	if err != nil {
		return nil, apperror.NewInternal("failed to create transaction", err)
	}

	u.TransactionCreated.SetPayload(transaction)
	u.EventDispatcher.Dispatch(u.TransactionCreated)

	return &CreateTransactionOutputDTO{
		ID: transaction.ID,
	}, nil
}
