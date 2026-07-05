package delete_transaction

import (
	"database/sql"
	"testing"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/apperror"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestDeleteTransactionUseCase(t *testing.T) {
	transaction := &entity.Transaction{ID: "tx-1"}

	transactionGateway := &TransactionGatewayMock{}
	transactionGateway.On("FindByID", transaction.ID).Return(transaction, nil)
	transactionGateway.On("Delete", transaction.ID).Return(nil)

	uc := NewDeleteTransactionUseCase(transactionGateway)
	output, err := uc.Execute(transaction.ID)

	assert.Nil(t, err)
	assert.Equal(t, transaction.ID, output.ID)
}

func TestDeleteTransactionUseCase_NotFound(t *testing.T) {
	transactionGateway := &TransactionGatewayMock{}
	transactionGateway.On("FindByID", "missing").Return((*entity.Transaction)(nil), sql.ErrNoRows)

	uc := NewDeleteTransactionUseCase(transactionGateway)
	output, err := uc.Execute("missing")

	assert.Nil(t, output)
	assert.Equal(t, apperror.CodeNotFound, apperror.CodeFromError(err))
}
