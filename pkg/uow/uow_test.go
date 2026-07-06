package uow

import (
	"database/sql"
	"testing"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/database"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/gateway"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestSQLUnitOfWork_Commit(t *testing.T) {
	db, accountFrom, accountTo := setupTransactionTestDB(t)
	unitOfWork := NewSQLUnitOfWork(db)

	transaction, err := entity.NewTransaction(accountFrom, accountTo, 100)
	assert.Nil(t, err)

	err = unitOfWork.Do(func(repos Repositories) error {
		if err := repos.Account.UpdateBalance(accountFrom); err != nil {
			return err
		}

		if err := repos.Account.UpdateBalance(accountTo); err != nil {
			return err
		}

		return repos.Transaction.Create(transaction)
	})
	assert.Nil(t, err)

	accountDB := database.NewAccountDB(db)
	from, err := accountDB.FindByID(accountFrom.ID)
	assert.Nil(t, err)
	assert.Equal(t, float64(900), from.Balance)

	to, err := accountDB.FindByID(accountTo.ID)
	assert.Nil(t, err)
	assert.Equal(t, float64(600), to.Balance)
}

func TestSQLUnitOfWork_Rollback(t *testing.T) {
	db, accountFrom, accountTo := setupTransactionTestDB(t)
	unitOfWork := NewSQLUnitOfWork(db)

	transaction, err := entity.NewTransaction(accountFrom, accountTo, 100)
	assert.Nil(t, err)

	err = unitOfWork.Do(func(repos Repositories) error {
		if err := repos.Account.UpdateBalance(accountFrom); err != nil {
			return err
		}

		if err := repos.Account.UpdateBalance(accountTo); err != nil {
			return err
		}

		_ = repos.Transaction.Create(transaction)
		return assert.AnError
	})
	assert.Error(t, err)

	accountDB := database.NewAccountDB(db)
	from, err := accountDB.FindByID(accountFrom.ID)
	assert.Nil(t, err)
	assert.Equal(t, float64(1000), from.Balance)

	transactionDB := database.NewTransactionDB(db)
	exists, err := transactionDB.Exists(accountFrom.ID, accountTo.ID, 100)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestPassThroughUnitOfWork(t *testing.T) {
	accountFrom := &entity.Account{ID: "from", Balance: 900}
	accountTo := &entity.Account{ID: "to", Balance: 1100}

	accountGateway := &accountGatewayStub{accountFrom: accountFrom, accountTo: accountTo}
	transactionGateway := &transactionGatewayStub{}

	unitOfWork := NewPassThroughUnitOfWork(accountGateway, transactionGateway)

	err := unitOfWork.Do(func(repos Repositories) error {
		assert.Equal(t, accountGateway, repos.Account)
		assert.Equal(t, transactionGateway, repos.Transaction)
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, accountGateway, unitOfWork.Account())
	assert.Equal(t, transactionGateway, unitOfWork.Transaction())
}

type accountGatewayStub struct {
	accountFrom *entity.Account
	accountTo   *entity.Account
}

func (s *accountGatewayStub) FindByID(id string) (*entity.Account, error) {
	if id == s.accountFrom.ID {
		return s.accountFrom, nil
	}

	return s.accountTo, nil
}

func (s *accountGatewayStub) FindByClientID(string) (*entity.Account, error) { return nil, nil }
func (s *accountGatewayStub) List() ([]*entity.Account, error)             { return nil, nil }
func (s *accountGatewayStub) Save(*entity.Account) error                     { return nil }
func (s *accountGatewayStub) UpdateBalance(*entity.Account) error            { return nil }
func (s *accountGatewayStub) Delete(string) error                            { return nil }
func (s *accountGatewayStub) DeleteAll() (int64, error)                      { return 0, nil }

type transactionGatewayStub struct{}

func (s *transactionGatewayStub) FindByID(string) (*entity.Transaction, error) { return nil, nil }
func (s *transactionGatewayStub) Create(*entity.Transaction) error             { return nil }
func (s *transactionGatewayStub) Exists(string, string, float64) (bool, error) {
	return false, nil
}
func (s *transactionGatewayStub) ExistsByAccountID(string) (bool, error) { return false, nil }
func (s *transactionGatewayStub) Delete(string) error                    { return nil }
func (s *transactionGatewayStub) DeleteAll() (int64, error)              { return 0, nil }

var (
	_ UnitOfWork = (*SQLUnitOfWork)(nil)
	_ UnitOfWork = (*PassThroughUnitOfWork)(nil)
	_ gateway.AccountGateway = (*accountGatewayStub)(nil)
	_ gateway.TransactionGateway = (*transactionGatewayStub)(nil)
)

func setupTransactionTestDB(t *testing.T) (*sql.DB, *entity.Account, *entity.Account) {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.Nil(t, err)

	t.Cleanup(func() {
		_ = db.Close()
	})

	_, err = db.Exec(`CREATE TABLE clients (
		id varchar(255), name varchar(255), email varchar(255),
		created_at datetime, updated_at datetime
	)`)
	assert.Nil(t, err)

	_, err = db.Exec(`CREATE TABLE accounts (
		id varchar(255), client_id varchar(255), balance float,
		created_at datetime, updated_at datetime
	)`)
	assert.Nil(t, err)

	_, err = db.Exec(`CREATE TABLE transactions (
		id varchar(255), account_from_id varchar(255), account_to_id varchar(255),
		amount float, created_at datetime
	)`)
	assert.Nil(t, err)

	client, _ := entity.NewClient("John", "j@j.com")
	_, err = db.Exec(
		"INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		client.ID, client.Name, client.Email, client.CreatedAt, client.UpdatedAt,
	)
	assert.Nil(t, err)

	accountFrom := entity.NewAccount(client)
	accountFrom.Balance = 1000

	accountTo := entity.NewAccount(client)
	accountTo.Balance = 500

	accountDB := database.NewAccountDB(db)
	assert.Nil(t, accountDB.Save(accountFrom))
	assert.Nil(t, accountDB.Save(accountTo))

	accountTo.Client = client
	accountFrom.Client = client

	return db, accountFrom, accountTo
}
