package delete_account

import (
	"database/sql"
	"testing"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
	return m.Called(account).Error(0)
}

func (m *AccountGatewayMock) UpdateBalance(account *entity.Account) error {
	return m.Called(account).Error(0)
}

func (m *AccountGatewayMock) Delete(id string) error {
	return m.Called(id).Error(0)
}

func (m *AccountGatewayMock) DeleteAll() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

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
	return m.Called(transaction).Error(0)
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

func TestDeleteAccountUseCase(t *testing.T) {
	client, _ := entity.NewClient("John", "j@j.com")
	account := entity.NewAccount(client)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", account.ID).Return(account, nil)
	accountGateway.On("Delete", account.ID).Return(nil)

	transactionGateway := &TransactionGatewayMock{}
	transactionGateway.On("ExistsByAccountID", account.ID).Return(false, nil)

	uc := NewDeleteAccountUseCase(accountGateway, transactionGateway)
	output, err := uc.Execute(account.ID)

	assert.Nil(t, err)
	assert.Equal(t, account.ID, output.ID)
}

func TestDeleteAccountUseCase_HasLinkedTransactions(t *testing.T) {
	client, _ := entity.NewClient("John", "j@j.com")
	account := entity.NewAccount(client)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", account.ID).Return(account, nil)

	transactionGateway := &TransactionGatewayMock{}
	transactionGateway.On("ExistsByAccountID", account.ID).Return(true, nil)

	uc := NewDeleteAccountUseCase(accountGateway, transactionGateway)
	output, err := uc.Execute(account.ID)

	assert.Nil(t, output)
	assert.Equal(t, apperror.CodeConflict, apperror.CodeFromError(err))
}

func TestDeleteAccountUseCase_NotFound(t *testing.T) {
	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", "missing").Return((*entity.Account)(nil), sql.ErrNoRows)

	uc := NewDeleteAccountUseCase(accountGateway, &TransactionGatewayMock{})
	output, err := uc.Execute("missing")

	assert.Nil(t, output)
	assert.Equal(t, apperror.CodeNotFound, apperror.CodeFromError(err))
}
