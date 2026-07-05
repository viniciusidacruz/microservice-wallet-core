package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID string
	Client *Client
	Balance float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccount(client *Client) *Account {
	if client == nil {
		return nil
	}

	account := &Account{
		ID: uuid.New().String(),
		Client: client,
		Balance: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return account
}

func (a *Account) Credit(amount float64) {
	a.Balance += amount
	a.UpdatedAt = time.Now()
}

func (a *Account) Debit(amount float64) {
	a.Balance -= amount
	a.UpdatedAt = time.Now()
}

func (a *Account) SetBalance(balance float64) error {
	if balance < 0 {
		return errors.New("balance must be greater than or equal to zero")
	}

	a.Balance = balance
	a.UpdatedAt = time.Now()

	return nil
}