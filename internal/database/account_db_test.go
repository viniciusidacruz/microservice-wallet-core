package database

import (
	"testing"

	"database/sql"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	DB        *sql.DB
	AccountDB *AccountDB
	Client    *entity.Client
}

func (s *AccountDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.DB = db
	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at datetime, updated_at datetime)")
	db.Exec("CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at datetime, updated_at datetime)")
	s.AccountDB = NewAccountDB(db)
	s.Client, _ = entity.NewClient("John Doe", "j@j.com")
}

func (s *AccountDBTestSuite) TearDownSuite() {
	defer s.DB.Close()
	s.DB.Exec("DROP TABLE accounts")
	s.DB.Exec("DROP TABLE clients")
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}

func (s *AccountDBTestSuite) TestSave() {
	account := entity.NewAccount(s.Client)
	err := s.AccountDB.Save(account)
	s.Nil(err)
}

func (s *AccountDBTestSuite) TestFindByID() {
	s.DB.Exec("INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", s.Client.ID, s.Client.Name, s.Client.Email, s.Client.CreatedAt, s.Client.UpdatedAt)

	account := entity.NewAccount(s.Client)
	err := s.AccountDB.Save(account)
	s.Nil(err)

	accountFound, err := s.AccountDB.FindByID(account.ID)
	s.Nil(err)
	s.Equal(account.ID, accountFound.ID)
	s.Equal(s.Client.ID, accountFound.Client.ID)
	s.Equal(account.Balance, accountFound.Balance)
	s.Equal(account.Client.ID, accountFound.Client.ID)
}

func (s *AccountDBTestSuite) TestUpdateBalance() {
	s.DB.Exec("INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", s.Client.ID, s.Client.Name, s.Client.Email, s.Client.CreatedAt, s.Client.UpdatedAt)

	account := entity.NewAccount(s.Client)
	err := s.AccountDB.Save(account)
	s.Nil(err)

	err = account.SetBalance(1500)
	s.Nil(err)

	err = s.AccountDB.UpdateBalance(account)
	s.Nil(err)

	accountFound, err := s.AccountDB.FindByID(account.ID)
	s.Nil(err)
	s.Equal(float64(1500), accountFound.Balance)
}

func (s *AccountDBTestSuite) TestList() {
	s.DB.Exec("INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", s.Client.ID, s.Client.Name, s.Client.Email, s.Client.CreatedAt, s.Client.UpdatedAt)

	account1 := entity.NewAccount(s.Client)
	s.Nil(s.AccountDB.Save(account1))

	client2, _ := entity.NewClient("Jane Doe", "j@j2.com")
	s.DB.Exec("INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", client2.ID, client2.Name, client2.Email, client2.CreatedAt, client2.UpdatedAt)

	account2 := entity.NewAccount(client2)
	s.Nil(s.AccountDB.Save(account2))

	accounts, err := s.AccountDB.List()
	s.Nil(err)
	s.Len(accounts, 2)
	s.Equal(account1.ID, accounts[0].ID)
	s.Equal(account2.ID, accounts[1].ID)
	s.Equal(s.Client.Name, accounts[0].Client.Name)
	s.Equal(client2.Name, accounts[1].Client.Name)
}

func (s *AccountDBTestSuite) TestFindByClientID() {
	s.DB.Exec("INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", s.Client.ID, s.Client.Name, s.Client.Email, s.Client.CreatedAt, s.Client.UpdatedAt)

	account := entity.NewAccount(s.Client)
	s.Nil(s.AccountDB.Save(account))

	accountFound, err := s.AccountDB.FindByClientID(s.Client.ID)
	s.Nil(err)
	s.Equal(account.ID, accountFound.ID)
}

func (s *AccountDBTestSuite) TestDelete() {
	s.DB.Exec("INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", s.Client.ID, s.Client.Name, s.Client.Email, s.Client.CreatedAt, s.Client.UpdatedAt)

	account := entity.NewAccount(s.Client)
	s.Nil(s.AccountDB.Save(account))
	s.Nil(s.AccountDB.Delete(account.ID))

	_, err := s.AccountDB.FindByID(account.ID)
	s.Error(err)
}
