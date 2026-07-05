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
	stmt, err := c.DB.Prepare("SELECT id, name, email, created_at, updated_at FROM clients WHERE id = ?")
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

func (c *ClientDB) GetByEmail(email string) (*entity.Client, error) {
	client := &entity.Client{}
	stmt, err := c.DB.Prepare("SELECT id, name, email, created_at, updated_at FROM clients WHERE email = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(email)

	if err := row.Scan(&client.ID, &client.Name, &client.Email, &client.CreatedAt, &client.UpdatedAt); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *ClientDB) List() ([]*entity.Client, error) {
	rows, err := c.DB.Query("SELECT id, name, email, created_at, updated_at FROM clients ORDER BY created_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := make([]*entity.Client, 0)

	for rows.Next() {
		client := &entity.Client{}
		if err := rows.Scan(&client.ID, &client.Name, &client.Email, &client.CreatedAt, &client.UpdatedAt); err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

func (c *ClientDB) Save(client *entity.Client) error {
	stmt, err := c.DB.Prepare("INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
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

func (c *ClientDB) Delete(id string) error {
	stmt, err := c.DB.Prepare("DELETE FROM clients WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (c *ClientDB) DeleteAll() (int64, error) {
	result, err := c.DB.Exec("DELETE FROM clients")
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
