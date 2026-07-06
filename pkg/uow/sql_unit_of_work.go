package uow

import (
	"database/sql"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/database"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
)

type SQLUnitOfWork struct {
	db          *sql.DB
	account     gateway.AccountGateway
	transaction gateway.TransactionGateway
}

func NewSQLUnitOfWork(db *sql.DB) *SQLUnitOfWork {
	return &SQLUnitOfWork{
		db:          db,
		account:     database.NewAccountDB(db),
		transaction: database.NewTransactionDB(db),
	}
}

func (u *SQLUnitOfWork) Account() gateway.AccountGateway {
	return u.account
}

func (u *SQLUnitOfWork) Transaction() gateway.TransactionGateway {
	return u.transaction
}

func (u *SQLUnitOfWork) Do(fn func(repos Repositories) error) error {
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	repos := Repositories{
		Account:     database.NewAccountDB(tx),
		Transaction: database.NewTransactionDB(tx),
	}

	if err := fn(repos); err != nil {
		return err
	}

	return tx.Commit()
}
