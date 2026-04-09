package tests

import (
	"testing"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	customerOne, _ := entity.NewClient("John", "test")
	accountOne := entity.NewAccount(customerOne)
	customerTwo, _ := entity.NewClient("John 2", "test2")
	accountTwo := entity.NewAccount(customerTwo)

	accountOne.Credit(1000)
	accountTwo.Credit(1000)

	transaction, err := entity.NewTransaction(accountOne, accountTwo, 100)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, 1100.0, accountTwo.Balance)
	assert.Equal(t, 900.0, accountOne.Balance)
}

func TestCreateTransactionWithInsufficientBalance(t *testing.T) {
	customerOne, _ := entity.NewClient("John", "test")
	accountOne := entity.NewAccount(customerOne)
	customerTwo, _ := entity.NewClient("John 2", "test2")
	accountTwo := entity.NewAccount(customerTwo)

	accountOne.Credit(1000)
	accountTwo.Credit(1000)

	transaction, err := entity.NewTransaction(accountOne, accountTwo, 2000)

	assert.NotNil(t, err)
	assert.Nil(t, transaction)
	assert.Equal(t, 1000.0, accountTwo.Balance)
	assert.Equal(t, 1000.0, accountOne.Balance)
	assert.Error(t, err, "Insufficient funds")
}