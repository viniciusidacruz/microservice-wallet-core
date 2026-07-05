package create_account

import (
	"database/sql"
	"testing"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ClientGatewayMock struct {
	mock.Mock
}

type AccountGatewayMock struct {
	mock.Mock
}

func (m *ClientGatewayMock) Save(client *entity.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *ClientGatewayMock) Get(id string) (*entity.Client, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Client), args.Error(1)
}

func (m *ClientGatewayMock) List() ([]*entity.Client, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*entity.Client), args.Error(1)
}

func (m *ClientGatewayMock) GetByEmail(email string) (*entity.Client, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Client), args.Error(1)
}

func (m *ClientGatewayMock) Delete(id string) error {
	return m.Called(id).Error(0)
}

func (m *ClientGatewayMock) DeleteAll() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
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

func TestCreateAccountUseCase(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	clientGateway := &ClientGatewayMock{}
	clientGateway.On("Get", client.ID).Return(client, nil)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByClientID", client.ID).Return((*entity.Account)(nil), sql.ErrNoRows)
	accountGateway.On("Save", mock.Anything).Return(nil)

	uc := NewCreateAccountUseCase(accountGateway, clientGateway)
	inputDTO := CreateAccountInputDTO{
		ClientID: client.ID,
	}
	output, err := uc.Execute(inputDTO)

	assert.Nil(t, err)
	assert.NotNil(t, output.ID)
	clientGateway.AssertExpectations(t)
	accountGateway.AssertExpectations(t)
	clientGateway.AssertNumberOfCalls(t, "Get", 1)
	accountGateway.AssertNumberOfCalls(t, "Save", 1)
}

func TestCreateAccountUseCase_EmptyClientID(t *testing.T) {
	uc := NewCreateAccountUseCase(&AccountGatewayMock{}, &ClientGatewayMock{})

	output, err := uc.Execute(CreateAccountInputDTO{})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeValidation, apperror.CodeFromError(err))
	assert.Equal(t, "client_id is required", err.Error())
}

func TestCreateAccountUseCase_ClientNotFound(t *testing.T) {
	clientGateway := &ClientGatewayMock{}
	clientGateway.On("Get", "missing-id").Return((*entity.Client)(nil), sql.ErrNoRows)

	uc := NewCreateAccountUseCase(&AccountGatewayMock{}, clientGateway)
	output, err := uc.Execute(CreateAccountInputDTO{ClientID: "missing-id"})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeNotFound, apperror.CodeFromError(err))
	assert.Equal(t, "client not found", err.Error())
}

func TestCreateAccountUseCase_SaveError(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	clientGateway := &ClientGatewayMock{}
	clientGateway.On("Get", client.ID).Return(client, nil)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByClientID", client.ID).Return((*entity.Account)(nil), sql.ErrNoRows)
	accountGateway.On("Save", mock.Anything).Return(sql.ErrConnDone)

	uc := NewCreateAccountUseCase(accountGateway, clientGateway)
	output, err := uc.Execute(CreateAccountInputDTO{ClientID: client.ID})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, apperror.CodeFromError(err))
	assert.Equal(t, "failed to save account", err.Error())
}

func TestCreateAccountUseCase_DuplicateAccount(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	existingAccount := entity.NewAccount(client)

	clientGateway := &ClientGatewayMock{}
	clientGateway.On("Get", client.ID).Return(client, nil)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByClientID", client.ID).Return(existingAccount, nil)

	uc := NewCreateAccountUseCase(accountGateway, clientGateway)
	output, err := uc.Execute(CreateAccountInputDTO{ClientID: client.ID})

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeConflict, apperror.CodeFromError(err))
	assert.Equal(t, "client already has an account", err.Error())
}
