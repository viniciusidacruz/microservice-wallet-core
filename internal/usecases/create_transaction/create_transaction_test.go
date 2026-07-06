package create_transaction

import (
	"database/sql"
	"testing"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/event"
	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/events"
	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/uow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) FindByID(id string) (*entity.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *TransactionGatewayMock) Exists(accountFromID, accountToID string, amount float64) (bool, error) {
	args := m.Called(accountFromID, accountToID, amount)
	return args.Bool(0), args.Error(1)
}

func (m *TransactionGatewayMock) ExistsByAccountID(accountID string) (bool, error) {
	args := m.Called(accountID)
	return args.Bool(0), args.Error(1)
}

func (m *TransactionGatewayMock) Delete(id string) error {
	return m.Called(id).Error(0)
}

func (m *TransactionGatewayMock) DeleteAll() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) FindByID(id string) (*entity.Account, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) FindByClientID(clientID string) (*entity.Account, error) {
	args := m.Called(clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) List() ([]*entity.Account, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) UpdateBalance(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) Delete(id string) error {
	return m.Called(id).Error(0)
}

func (m *AccountGatewayMock) DeleteAll() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateTransactionUseCase(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "j@j.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("John Doe 2", "j@j2.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", account1.ID).Return(account1, nil)
	accountGateway.On("FindByID", account2.ID).Return(account2, nil)
	accountGateway.On("UpdateBalance", mock.MatchedBy(func(account *entity.Account) bool {
		return account.ID == account1.ID && account.Balance == 900
	})).Return(nil)
	accountGateway.On("UpdateBalance", mock.MatchedBy(func(account *entity.Account) bool {
		return account.ID == account2.ID && account.Balance == 1100
	})).Return(nil)

	transactionGateway := &TransactionGatewayMock{}
	transactionGateway.On("Exists", account1.ID, account2.ID, float64(100)).Return(false, nil)
	transactionGateway.On("Create", mock.Anything).Return(nil)

	inputDTO := CreateTransactionInputDTO{
		AccountFromID: account1.ID,
		AccountToID:   account2.ID,
		Amount:        100,
	}

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	unitOfWork := uow.NewPassThroughUnitOfWork(accountGateway, transactionGateway)
	uc := NewCreateTransactionUseCase(unitOfWork, eventDispatcher, transactionCreatedEvent)
	output, err := uc.Execute(inputDTO)

	assert.Nil(t, err)
	assert.NotNil(t, output.ID)
	accountGateway.AssertExpectations(t)
	transactionGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "FindByID", 2)
	accountGateway.AssertNumberOfCalls(t, "UpdateBalance", 2)
	transactionGateway.AssertNumberOfCalls(t, "Create", 1)
}

func TestCreateTransactionUseCase_ValidationErrors(t *testing.T) {
	uc := NewCreateTransactionUseCase(nil, events.NewEventDispatcher(), event.NewTransactionCreated())

	tests := []struct {
		name    string
		input   CreateTransactionInputDTO
		message string
	}{
		{
			name:    "missing account from",
			input:   CreateTransactionInputDTO{AccountToID: "acc-2", Amount: 10},
			message: "account_from_id is required",
		},
		{
			name:    "missing account to",
			input:   CreateTransactionInputDTO{AccountFromID: "acc-1", Amount: 10},
			message: "account_to_id is required",
		},
		{
			name:    "same accounts",
			input:   CreateTransactionInputDTO{AccountFromID: "acc-1", AccountToID: "acc-1", Amount: 10},
			message: "account_from_id and account_to_id must be different",
		},
		{
			name:    "invalid amount",
			input:   CreateTransactionInputDTO{AccountFromID: "acc-1", AccountToID: "acc-2", Amount: 0},
			message: "amount must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := uc.Execute(tt.input)

			assert.Nil(t, output)
			assert.Error(t, err)
			assert.Equal(t, apperror.CodeValidation, apperror.CodeFromError(err))
			assert.Equal(t, tt.message, err.Error())
		})
	}
}

func TestCreateTransactionUseCase_AccountFromNotFound(t *testing.T) {
	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", "missing-from").Return((*entity.Account)(nil), sql.ErrNoRows)

	unitOfWork := uow.NewPassThroughUnitOfWork(accountGateway, &TransactionGatewayMock{})
	uc := NewCreateTransactionUseCase(unitOfWork, events.NewEventDispatcher(), event.NewTransactionCreated())
	output, err := uc.Execute(CreateTransactionInputDTO{
		AccountFromID: "missing-from",
		AccountToID:   "acc-2",
		Amount:        10,
	})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeNotFound, apperror.CodeFromError(err))
	assert.Equal(t, "account from not found", err.Error())
}

func TestCreateTransactionUseCase_InsufficientFunds(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "j@j.com")
	account1 := entity.NewAccount(client1)

	client2, _ := entity.NewClient("Jane Doe", "j@j2.com")
	account2 := entity.NewAccount(client2)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", account1.ID).Return(account1, nil)
	accountGateway.On("FindByID", account2.ID).Return(account2, nil)

	unitOfWork := uow.NewPassThroughUnitOfWork(accountGateway, &TransactionGatewayMock{})
	uc := NewCreateTransactionUseCase(unitOfWork, events.NewEventDispatcher(), event.NewTransactionCreated())
	output, err := uc.Execute(CreateTransactionInputDTO{
		AccountFromID: account1.ID,
		AccountToID:   account2.ID,
		Amount:        100,
	})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeValidation, apperror.CodeFromError(err))
	assert.Equal(t, "Insufficient funds", err.Error())
}

func TestCreateTransactionUseCase_DuplicateTransaction(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "j@j.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("Jane Doe", "j@j2.com")
	account2 := entity.NewAccount(client2)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", account1.ID).Return(account1, nil)
	accountGateway.On("FindByID", account2.ID).Return(account2, nil)

	transactionGateway := &TransactionGatewayMock{}
	transactionGateway.On("Exists", account1.ID, account2.ID, float64(100)).Return(true, nil)

	uc := NewCreateTransactionUseCase(uow.NewPassThroughUnitOfWork(accountGateway, transactionGateway), events.NewEventDispatcher(), event.NewTransactionCreated())
	output, err := uc.Execute(CreateTransactionInputDTO{
		AccountFromID: account1.ID,
		AccountToID:   account2.ID,
		Amount:        100,
	})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeConflict, apperror.CodeFromError(err))
	assert.Equal(t, "transaction already exists", err.Error())
}
