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

func (s *ClientDBTestSuite) TestList() {
	client1, _ := entity.NewClient("John Doe", "j@j.com")
	client2, _ := entity.NewClient("Jane Doe", "j@j2.com")

	s.Nil(s.ClientDB.Save(client1))
	s.Nil(s.ClientDB.Save(client2))

	clients, err := s.ClientDB.List()
	s.Nil(err)
	s.Len(clients, 2)
	s.Equal(client1.ID, clients[0].ID)
	s.Equal(client2.ID, clients[1].ID)
}

func (s *ClientDBTestSuite) TestGetByEmail() {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	s.Nil(s.ClientDB.Save(client))

	clientFound, err := s.ClientDB.GetByEmail(client.Email)
	s.Nil(err)
	s.Equal(client.ID, clientFound.ID)
}

func (s *ClientDBTestSuite) TestDelete() {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	s.Nil(s.ClientDB.Save(client))
	s.Nil(s.ClientDB.Delete(client.ID))

	_, err := s.ClientDB.Get(client.ID)
	s.Error(err)
}

func (s *ClientDBTestSuite) TestDeleteAll() {
	client1, _ := entity.NewClient("John Doe", "j@j.com")
	client2, _ := entity.NewClient("Jane Doe", "j@j2.com")
	s.Nil(s.ClientDB.Save(client1))
	s.Nil(s.ClientDB.Save(client2))

	count, err := s.ClientDB.DeleteAll()
	s.Nil(err)
	s.Equal(int64(2), count)
}
