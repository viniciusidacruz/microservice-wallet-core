package delete_all_accounts

import (
	"testing"

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

func TestDeleteAllAccountsUseCase(t *testing.T) {
	transactionGateway := &TransactionGatewayMock{}
	transactionGateway.On("DeleteAll").Return(int64(3), nil)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("DeleteAll").Return(int64(2), nil)

	uc := NewDeleteAllAccountsUseCase(accountGateway, transactionGateway)
	output, err := uc.Execute()

	assert.Nil(t, err)
	assert.Equal(t, int64(2), output.DeletedCount)
}
