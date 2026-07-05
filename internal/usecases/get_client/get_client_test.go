package get_client

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

func (m *ClientGatewayMock) Delete(id string) error {
	return m.Called(id).Error(0)
}

func (m *ClientGatewayMock) DeleteAll() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func TestGetClientUseCase(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	clientGateway := &ClientGatewayMock{}
	clientGateway.On("Get", client.ID).Return(client, nil)

	uc := NewGetClientUseCase(clientGateway)
	output, err := uc.Execute(client.ID)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, client.ID, output.ID)
	assert.Equal(t, client.Name, output.Name)
	assert.Equal(t, client.Email, output.Email)
	clientGateway.AssertExpectations(t)
}

func TestGetClientUseCase_EmptyID(t *testing.T) {
	uc := NewGetClientUseCase(&ClientGatewayMock{})
	output, err := uc.Execute("")

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeValidation, apperror.CodeFromError(err))
	assert.Equal(t, "client id is required", err.Error())
}

func TestGetClientUseCase_NotFound(t *testing.T) {
	clientGateway := &ClientGatewayMock{}
	clientGateway.On("Get", "missing-id").Return((*entity.Client)(nil), sql.ErrNoRows)

	uc := NewGetClientUseCase(clientGateway)
	output, err := uc.Execute("missing-id")

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeNotFound, apperror.CodeFromError(err))
	assert.Equal(t, "client not found", err.Error())
}

func TestGetClientUseCase_InternalError(t *testing.T) {
	clientGateway := &ClientGatewayMock{}
	clientGateway.On("Get", "client-id").Return((*entity.Client)(nil), sql.ErrConnDone)

	uc := NewGetClientUseCase(clientGateway)
	output, err := uc.Execute("client-id")

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, apperror.CodeFromError(err))
	assert.Equal(t, "failed to get client", err.Error())
}
