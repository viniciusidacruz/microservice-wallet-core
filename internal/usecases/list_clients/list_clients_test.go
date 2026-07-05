package list_clients

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

func TestListClientsUseCase(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "j@j.com")
	client2, _ := entity.NewClient("Jane Doe", "j@j2.com")

	clientGateway := &ClientGatewayMock{}
	clientGateway.On("List").Return([]*entity.Client{client1, client2}, nil)

	uc := NewListClientsUseCase(clientGateway)
	output, err := uc.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Clients, 2)
	assert.Equal(t, client1.ID, output.Clients[0].ID)
	assert.Equal(t, client2.ID, output.Clients[1].ID)
	clientGateway.AssertExpectations(t)
}

func TestListClientsUseCase_EmptyList(t *testing.T) {
	clientGateway := &ClientGatewayMock{}
	clientGateway.On("List").Return([]*entity.Client{}, nil)

	uc := NewListClientsUseCase(clientGateway)
	output, err := uc.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Empty(t, output.Clients)
}

func TestListClientsUseCase_InternalError(t *testing.T) {
	clientGateway := &ClientGatewayMock{}
	clientGateway.On("List").Return(([]*entity.Client)(nil), sql.ErrConnDone)

	uc := NewListClientsUseCase(clientGateway)
	output, err := uc.Execute()

	assert.Nil(t, output)
	assert.Error(t, err)
	assert.Equal(t, apperror.CodeInternal, apperror.CodeFromError(err))
	assert.Equal(t, "failed to list clients", err.Error())
}
