package database

import (
	"database/sql"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
)

type ClientDB struct {
	DB *sql.DB
}

func NewClientDB(db *sql.DB) *ClientDB {
	return &ClientDB{
		DB: db,
	}
}

func (c *ClientDB) Get(id string) (*entity.Client, error) {
	client := &entity.Client{}
	stmt, err := c.DB.Prepare("SELECT id, name, email, created_at, updated_at FROM clients WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	if err := row.Scan(&client.ID, &client.Name, &client.Email, &client.CreatedAt, &client.UpdatedAt); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *ClientDB) Save(client *entity.Client) error {
	stmt, err := c.DB.Prepare("INSERT INTO clients (id, name, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(client.ID, client.Name, client.Email, client.CreatedAt, client.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
