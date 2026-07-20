package create_transaction

import (
	"database/sql"
	"errors"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/events"
	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/uow"
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
	UnitOfWork         uow.UnitOfWork
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	BalanceUpdated     events.EventInterface
}

type BalanceUpdatedOutputDTO struct {
	AccountFromID      string  `json:"account_from_id"`
	AccountToID        string  `json:"account_to_id"`
	BalanceAccountFrom float64 `json:"balance_account_from"`
	BalanceAccountTo   float64 `json:"balance_account_to"`
}

func NewCreateTransactionUseCase(
	unitOfWork uow.UnitOfWork,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreatedEvent events.EventInterface,
	balanceUpdatedEvent events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		UnitOfWork:         unitOfWork,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreatedEvent,
		BalanceUpdated:     balanceUpdatedEvent,
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

	accountFrom, err := u.UnitOfWork.Account().FindByID(input.AccountFromID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("account from not found")
		}

		return nil, apperror.NewInternal("failed to get account from", err)
	}

	accountTo, err := u.UnitOfWork.Account().FindByID(input.AccountToID)
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

	exists, err := u.UnitOfWork.Transaction().Exists(input.AccountFromID, input.AccountToID, input.Amount)
	if err != nil {
		return nil, apperror.NewInternal("failed to check transaction duplicate", err)
	}

	if exists {
		return nil, apperror.NewConflict("transaction already exists")
	}

	err = u.UnitOfWork.Do(func(repos uow.Repositories) error {
		if err := repos.Account.UpdateBalance(accountFrom); err != nil {
			return apperror.NewInternal("failed to update account from balance", err)
		}

		if err := repos.Account.UpdateBalance(accountTo); err != nil {
			return apperror.NewInternal("failed to update account to balance", err)
		}

		if err := repos.Transaction.Create(transaction); err != nil {
			return apperror.NewInternal("failed to create transaction", err)
		}

		return nil
	})
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			return nil, appErr
		}

		return nil, apperror.NewInternal("failed to commit transaction", err)
	}

	u.TransactionCreated.SetPayload(transaction)
	u.EventDispatcher.Dispatch(u.TransactionCreated)

	u.BalanceUpdated.SetPayload(BalanceUpdatedOutputDTO{
		AccountFromID:      accountFrom.ID,
		AccountToID:        accountTo.ID,
		BalanceAccountFrom: accountFrom.Balance,
		BalanceAccountTo:   accountTo.Balance,
	})
	u.EventDispatcher.Dispatch(u.BalanceUpdated)

	return &CreateTransactionOutputDTO{
		ID: transaction.ID,
	}, nil
}
