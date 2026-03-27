package repository

import (
	"database/sql"
	"errors"

	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
)

type DBOrTx interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func GetAllAccounts(conn DBOrTx) ([]models.Account, error) {
	var accounts []models.Account
	err := conn.Select(&accounts, "SELECT id,balance,owner FROM accounts")
	return accounts, err
}

func GetAccountsByOwner(conn DBOrTx, ownerName string) ([]models.Account, error) {
	var accounts []models.Account
	err := conn.Select(&accounts, "SELECT id, balance, owner FROM accounts WHERE owner ILIKE $1", "%"+ownerName+"%")
	return accounts, err
}

func GetAccountByID(conn DBOrTx, id int) (models.Account, error) {
	var account models.Account
	err := conn.Get(&account, "SELECT id,balance,owner FROM accounts WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Account{}, errors.New("счет с таким ID не существует")
		}
		return models.Account{}, err
	}

	return account, nil
}

func TopUpAccount(conn DBOrTx, accountID int, amount float64) error {
	_, err := conn.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, accountID)
	return err
}

func WithdrawAccount(conn DBOrTx, accountID int, amount float64) error {
	_, err := conn.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, accountID)
	return err
}
