package delete_client

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

func (m *ClientGatewayMock) Get(id string) (*entity.Client, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Client), args.Error(1)
}

func (m *ClientGatewayMock) GetByEmail(email string) (*entity.Client, error) {
	args := m.Called(email)
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

func (m *ClientGatewayMock) Save(client *entity.Client) error {
	return m.Called(client).Error(0)
}

func (m *ClientGatewayMock) Delete(id string) error {
	return m.Called(id).Error(0)
}

func (m *ClientGatewayMock) DeleteAll() (int64, error) {
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

func TestDeleteClientUseCase(t *testing.T) {
	client, _ := entity.NewClient("John", "j@j.com")
	clientGateway := &ClientGatewayMock{}
	clientGateway.On("Get", client.ID).Return(client, nil)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByClientID", client.ID).Return((*entity.Account)(nil), sql.ErrNoRows)
	clientGateway.On("Delete", client.ID).Return(nil)

	uc := NewDeleteClientUseCase(clientGateway, accountGateway)
	output, err := uc.Execute(client.ID)

	assert.Nil(t, err)
	assert.Equal(t, client.ID, output.ID)
}

func TestDeleteClientUseCase_NotFound(t *testing.T) {
	clientGateway := &ClientGatewayMock{}
	clientGateway.On("Get", "missing").Return((*entity.Client)(nil), sql.ErrNoRows)

	uc := NewDeleteClientUseCase(clientGateway, &AccountGatewayMock{})
	output, err := uc.Execute("missing")

	assert.Nil(t, output)
	assert.Equal(t, apperror.CodeNotFound, apperror.CodeFromError(err))
}

func TestDeleteClientUseCase_HasLinkedAccounts(t *testing.T) {
	client, _ := entity.NewClient("John", "j@j.com")
	account := entity.NewAccount(client)

	clientGateway := &ClientGatewayMock{}
	clientGateway.On("Get", client.ID).Return(client, nil)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByClientID", client.ID).Return(account, nil)

	uc := NewDeleteClientUseCase(clientGateway, accountGateway)
	output, err := uc.Execute(client.ID)

	assert.Nil(t, output)
	assert.Equal(t, apperror.CodeConflict, apperror.CodeFromError(err))
	assert.Equal(t, "client has linked accounts", err.Error())
}
