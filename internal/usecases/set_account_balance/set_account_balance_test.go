package set_account_balance

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

func TestSetAccountBalanceUseCase(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	account := entity.NewAccount(client)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", account.ID).Return(account, nil)
	accountGateway.On("UpdateBalance", mock.MatchedBy(func(updated *entity.Account) bool {
		return updated.ID == account.ID && updated.Balance == 1000
	})).Return(nil)

	uc := NewSetAccountBalanceUseCase(accountGateway)
	output, err := uc.Execute(SetAccountBalanceInputDTO{
		AccountID: account.ID,
		Balance:   1000,
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, account.ID, output.AccountID)
	assert.Equal(t, float64(1000), output.Balance)
	accountGateway.AssertExpectations(t)
}

func TestSetAccountBalanceUseCase_EmptyAccountID(t *testing.T) {
	uc := NewSetAccountBalanceUseCase(&AccountGatewayMock{})

	output, err := uc.Execute(SetAccountBalanceInputDTO{Balance: 100})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeValidation, apperror.CodeFromError(err))
	assert.Equal(t, "account_id is required", err.Error())
}

func TestSetAccountBalanceUseCase_NegativeBalance(t *testing.T) {
	uc := NewSetAccountBalanceUseCase(&AccountGatewayMock{})

	output, err := uc.Execute(SetAccountBalanceInputDTO{
		AccountID: "acc-1",
		Balance:   -10,
	})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeValidation, apperror.CodeFromError(err))
	assert.Equal(t, "balance must be greater than or equal to zero", err.Error())
}

func TestSetAccountBalanceUseCase_AccountNotFound(t *testing.T) {
	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", "missing-id").Return((*entity.Account)(nil), sql.ErrNoRows)

	uc := NewSetAccountBalanceUseCase(accountGateway)
	output, err := uc.Execute(SetAccountBalanceInputDTO{
		AccountID: "missing-id",
		Balance:   100,
	})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeNotFound, apperror.CodeFromError(err))
	assert.Equal(t, "account not found", err.Error())
}

func TestSetAccountBalanceUseCase_UpdateError(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	account := entity.NewAccount(client)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", account.ID).Return(account, nil)
	accountGateway.On("UpdateBalance", mock.Anything).Return(sql.ErrConnDone)

	uc := NewSetAccountBalanceUseCase(accountGateway)
	output, err := uc.Execute(SetAccountBalanceInputDTO{
		AccountID: account.ID,
		Balance:   500,
	})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, apperror.CodeFromError(err))
	assert.Equal(t, "failed to update account balance", err.Error())
}
