package tests

import (
	"testing"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewClient(t *testing.T) {
	client, err := entity.NewClient("John Doe", "j@j.com")

	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "John Doe", client.Name)
	assert.Equal(t, "j@j.com", client.Email)
}

func TestCreateNewClientWhenArgsAreInvalid(t *testing.T) {
	client, err := entity.NewClient("", "")

	assert.Nil(t, client)
	assert.NotNil(t, err)
}

func TestUpdateClient(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	err := client.Update("John Doe Updated", "j@jupdated.com")

	assert.Nil(t, err)
	assert.Equal(t, "John Doe Updated", client.Name)
	assert.Equal(t, "j@jupdated.com", client.Email)

}

func TestUpdateClientWhenArgsAreInvalid(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	err := client.Update("", "j@jupdated.com")

	assert.Error(t, err, "Name is required")
}

func TestAddAccountToClient(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	account := entity.NewAccount(client)

	err := client.AddAccount(account)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(client.Accounts))
}