package createtransaction

import (
	"testing"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *AccountGatewayMock) FindByID(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

func TestCreateTransactionUseCase(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "j@j.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("John Doe 2", "j@j2.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	accountGateway := &AccountGatewayMock{}
	accountGateway.On("FindByID", account1.ID).Return(account1, nil)
	accountGateway.On("FindByID", account2.ID).Return(account2, nil)

	transactionGateway := &TransactionGatewayMock{}
	transactionGateway.On("Create", mock.Anything).Return(nil)

	inputDTO := CreateTransactionInputDTO{
		AccountFromID: account1.ID,
		AccountToID:   account2.ID,
		Amount:        100,
	}
	uc := NewCreateTransactionUseCase(transactionGateway, accountGateway)
	output, err := uc.Execute(inputDTO)

	assert.Nil(t, err)
	assert.NotNil(t, output.ID)
	accountGateway.AssertExpectations(t)
	transactionGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "FindByID", 2)
	transactionGateway.AssertNumberOfCalls(t, "Create", 1)
}
