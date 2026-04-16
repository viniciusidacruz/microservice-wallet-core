package createaccount

import (
	"testing"

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

	return args.Get(0).(*entity.Client), args.Error(1)
}

func (m *AccountGatewayMock) FindByID(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func TestCreateAccountUseCase(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	clientGateway := &ClientGatewayMock{}
	clientGateway.On("Get", client.ID).Return(client, nil)

	accountGateway := &AccountGatewayMock{}
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