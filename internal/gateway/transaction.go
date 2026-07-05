package gateway

import "github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"

type TransactionGateway interface {
	FindByID(id string) (*entity.Transaction, error)
	Create(transaction *entity.Transaction) error
	Exists(accountFromID, accountToID string, amount float64) (bool, error)
	ExistsByAccountID(accountID string) (bool, error)
	Delete(id string) error
	DeleteAll() (int64, error)
}