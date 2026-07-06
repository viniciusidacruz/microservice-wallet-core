package database

import (
	"database/sql"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
)

type AccountDB struct {
	DB DBTX
}

func NewAccountDB(db DBTX) *AccountDB {
	return &AccountDB{
		DB: db,
	}
}

func (a *AccountDB) FindByID(id string) (*entity.Account, error) {
	var account entity.Account
	var client entity.Client
	account.Client = &client

	stmt, err := a.DB.Prepare("SELECT a.id, a.client_id, a.balance, a.created_at, a.updated_at, c.id, c.name, c.email, c.created_at, c.updated_at FROM accounts a INNER JOIN clients c ON a.client_id = c.id WHERE a.id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	if err := row.Scan(&account.ID, &account.Client.ID, &account.Balance, &account.CreatedAt, &account.UpdatedAt, &client.ID, &client.Name, &client.Email, &client.CreatedAt, &client.UpdatedAt); err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountDB) FindByClientID(clientID string) (*entity.Account, error) {
	var account entity.Account
	var client entity.Client
	account.Client = &client

	stmt, err := a.DB.Prepare("SELECT a.id, a.client_id, a.balance, a.created_at, a.updated_at, c.id, c.name, c.email, c.created_at, c.updated_at FROM accounts a INNER JOIN clients c ON a.client_id = c.id WHERE a.client_id = ? LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(clientID)

	if err := row.Scan(&account.ID, &account.Client.ID, &account.Balance, &account.CreatedAt, &account.UpdatedAt, &client.ID, &client.Name, &client.Email, &client.CreatedAt, &client.UpdatedAt); err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountDB) List() ([]*entity.Account, error) {
	rows, err := a.DB.Query("SELECT a.id, a.client_id, a.balance, a.created_at, a.updated_at, c.id, c.name, c.email, c.created_at, c.updated_at FROM accounts a INNER JOIN clients c ON a.client_id = c.id ORDER BY a.created_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*entity.Account, 0)

	for rows.Next() {
		var account entity.Account
		var client entity.Client
		account.Client = &client

		if err := rows.Scan(&account.ID, &account.Client.ID, &account.Balance, &account.CreatedAt, &account.UpdatedAt, &client.ID, &client.Name, &client.Email, &client.CreatedAt, &client.UpdatedAt); err != nil {
			return nil, err
		}

		accounts = append(accounts, &account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (a *AccountDB) Save(account *entity.Account) error {
	stmt, err := a.DB.Prepare("INSERT INTO accounts (id, client_id, balance, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.ID, account.Client.ID, account.Balance, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountDB) UpdateBalance(account *entity.Account) error {
	stmt, err := a.DB.Prepare("UPDATE accounts SET balance = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(account.Balance, account.UpdatedAt, account.ID)
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

func (a *AccountDB) Delete(id string) error {
	stmt, err := a.DB.Prepare("DELETE FROM accounts WHERE id = ?")
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

func (a *AccountDB) DeleteAll() (int64, error) {
	result, err := a.DB.Exec("DELETE FROM accounts")
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
