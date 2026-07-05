package get_account

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

func TestGetAccountUseCase(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	account := entity.NewAccount(client)
	account.Credit(500)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", account.ID).Return(account, nil)

	uc := NewGetAccountUseCase(accountGateway)
	output, err := uc.Execute(account.ID)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, account.ID, output.ID)
	assert.Equal(t, client.ID, output.ClientID)
	assert.Equal(t, float64(500), output.Balance)
	assert.Equal(t, client.Name, output.Client.Name)
	accountGateway.AssertExpectations(t)
}

func TestGetAccountUseCase_EmptyID(t *testing.T) {
	uc := NewGetAccountUseCase(&AccountGatewayMock{})
	output, err := uc.Execute("")

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeValidation, apperror.CodeFromError(err))
	assert.Equal(t, "account id is required", err.Error())
}

func TestGetAccountUseCase_NotFound(t *testing.T) {
	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", "missing-id").Return((*entity.Account)(nil), sql.ErrNoRows)

	uc := NewGetAccountUseCase(accountGateway)
	output, err := uc.Execute("missing-id")

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeNotFound, apperror.CodeFromError(err))
	assert.Equal(t, "account not found", err.Error())
}

func TestGetAccountUseCase_InternalError(t *testing.T) {
	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", "account-id").Return((*entity.Account)(nil), sql.ErrConnDone)

	uc := NewGetAccountUseCase(accountGateway)
	output, err := uc.Execute("account-id")

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, apperror.CodeFromError(err))
	assert.Equal(t, "failed to get account", err.Error())
}
