package list_accounts

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

func TestListAccountsUseCase(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "j@j.com")
	account1 := entity.NewAccount(client1)

	client2, _ := entity.NewClient("Jane Doe", "j@j2.com")
	account2 := entity.NewAccount(client2)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("List").Return([]*entity.Account{account1, account2}, nil)

	uc := NewListAccountsUseCase(accountGateway)
	output, err := uc.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Accounts, 2)
	assert.Equal(t, account1.ID, output.Accounts[0].ID)
	assert.Equal(t, account2.ID, output.Accounts[1].ID)
	accountGateway.AssertExpectations(t)
}

func TestListAccountsUseCase_EmptyList(t *testing.T) {
	accountGateway := &AccountGatewayMock{}
	accountGateway.On("List").Return([]*entity.Account{}, nil)

	uc := NewListAccountsUseCase(accountGateway)
	output, err := uc.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Empty(t, output.Accounts)
}

func TestListAccountsUseCase_InternalError(t *testing.T) {
	accountGateway := &AccountGatewayMock{}
	accountGateway.On("List").Return(([]*entity.Account)(nil), sql.ErrConnDone)

	uc := NewListAccountsUseCase(accountGateway)
	output, err := uc.Execute()

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, apperror.CodeFromError(err))
	assert.Equal(t, "failed to list accounts", err.Error())
}
