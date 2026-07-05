package tests

import (
	"testing"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	account := entity.NewAccount(client)
	assert.NotNil(t, account)
	assert.Equal(t, client.ID, account.Client.ID)
}

func TestCreateNewAccountWhenArgsAreInvalid(t *testing.T) {
	account := entity.NewAccount(nil)

	assert.Nil(t, account)
}

func TestCreditAccount(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	account := entity.NewAccount(client)
	account.Credit(100)
	assert.Equal(t, float64(100), account.Balance)
}

func TestDebitAccount(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	account := entity.NewAccount(client)
	account.Credit(100)
	account.Debit(50)
	assert.Equal(t, float64(50), account.Balance)
}

func TestSetAccountBalance(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	account := entity.NewAccount(client)

	err := account.SetBalance(250)
	assert.Nil(t, err)
	assert.Equal(t, float64(250), account.Balance)
}

func TestSetAccountBalanceWhenBalanceIsNegative(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	account := entity.NewAccount(client)

	err := account.SetBalance(-1)
	assert.Error(t, err)
	assert.Equal(t, "balance must be greater than or equal to zero", err.Error())
}