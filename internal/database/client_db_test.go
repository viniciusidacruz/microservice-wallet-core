package database

import (
	"database/sql"
	"testing"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type ClientDBTestSuite struct {
	suite.Suite
	DB       *sql.DB
	ClientDB *ClientDB
}

func (s *ClientDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.DB = db
	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at datetime, updated_at datetime)")
	s.ClientDB = NewClientDB(db)
}

func (s *ClientDBTestSuite) TearDownSuite() {
	defer s.DB.Close()
	s.DB.Exec("DROP TABLE clients")
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuite))
}

func (s *ClientDBTestSuite) TestGet() {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	s.ClientDB.Save(client)

	clientFound, err := s.ClientDB.Get(client.ID)
	s.Nil(err)
	s.Equal(client.ID, clientFound.ID)
	s.Equal(client.Name, clientFound.Name)
	s.Equal(client.Email, clientFound.Email)
}

func (s *ClientDBTestSuite) TestSave() {
	client := &entity.Client{
		ID:    "1",
		Name:  "John Doe",
		Email: "j@j.com",
	}
	err := s.ClientDB.Save(client)
	s.Nil(err)

	clientFound, err := s.ClientDB.Get(client.ID)
	s.Nil(err)
	s.Equal(client.ID, clientFound.ID)
	s.Equal(client.Name, clientFound.Name)
	s.Equal(client.Email, clientFound.Email)
}
