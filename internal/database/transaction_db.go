package database

import (
	"database/sql"

	"github.com.br/viniciusidacruz/microservice-wallet-core/internal/entity"
)

type TransactionDB struct {
	DB DBTX
}

func NewTransactionDB(db DBTX) *TransactionDB {
	return &TransactionDB{
		DB: db,
	}
}

func (t *TransactionDB) FindByID(id string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{
		AccountFrom: &entity.Account{},
		AccountTo:   &entity.Account{},
	}

	stmt, err := t.DB.Prepare("SELECT id, account_from_id, account_to_id, amount, created_at FROM transactions WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	if err := row.Scan(&transaction.ID, &transaction.AccountFrom.ID, &transaction.AccountTo.ID, &transaction.Amount, &transaction.CreatedAt); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionDB) Create(transaction *entity.Transaction) error {
	stmt, err := t.DB.Prepare("INSERT INTO transactions (id, account_from_id, account_to_id, amount, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(transaction.ID, transaction.AccountFrom.ID, transaction.AccountTo.ID, transaction.Amount, transaction.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionDB) Exists(accountFromID, accountToID string, amount float64) (bool, error) {
	stmt, err := t.DB.Prepare("SELECT COUNT(1) FROM transactions WHERE account_from_id = ? AND account_to_id = ? AND amount = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(accountFromID, accountToID, amount).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (t *TransactionDB) ExistsByAccountID(accountID string) (bool, error) {
	stmt, err := t.DB.Prepare("SELECT COUNT(1) FROM transactions WHERE account_from_id = ? OR account_to_id = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(accountID, accountID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (t *TransactionDB) Delete(id string) error {
	stmt, err := t.DB.Prepare("DELETE FROM transactions WHERE id = ?")
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

func (t *TransactionDB) DeleteAll() (int64, error) {
	result, err := t.DB.Exec("DELETE FROM transactions")
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
