package database

import (
	"database/sql"
	"testing"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	DB            *sql.DB
	TransactionDB *TransactionDB
	client        *entity.Client
	client2       *entity.Client
	accountFrom   *entity.Account
	accountTo     *entity.Account
}

func (s *TransactionDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.DB = db
	db.Exec("CREATE TABLE transactions (id varchar(255), account_from_id varchar(255), account_to_id varchar(255), amount float, created_at datetime)")
	db.Exec("CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at datetime, updated_at datetime)")
	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at datetime, updated_at datetime)")
	client, err := entity.NewClient("John Doe", "j@j.com")
	s.Nil(err)
	s.client = client

	client2, err := entity.NewClient("John Doe 2", "j@j2.com")
	s.Nil(err)
	s.client2 = client2

	accountFrom := entity.NewAccount(client)
	accountFrom.Balance = 1000
	s.accountFrom = accountFrom
	accountTo := entity.NewAccount(client2)
	accountTo.Balance = 1000
	s.accountTo = accountTo

	s.TransactionDB = NewTransactionDB(db)
}

func (s *TransactionDBTestSuite) TearDownSuite() {
	defer s.DB.Close()
	s.DB.Exec("DROP TABLE transactions")
	s.DB.Exec("DROP TABLE accounts")
	s.DB.Exec("DROP TABLE clients")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(s.accountFrom, s.accountTo, 100)
	s.Nil(err)
	err = s.TransactionDB.Create(transaction)
	s.Nil(err)
	s.Equal(float64(900), s.accountFrom.Balance)
	s.Equal(float64(1100), s.accountTo.Balance)
}
