package repository

import (
	"database/sql"
	"errors"
	"log"

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
	if err != nil {
		log.Printf("[ERROR] GetAllAccounts query failed: %v", err)
		return nil, err
	}
	return accounts, err
}

func GetAccountsByOwner(conn DBOrTx, ownerName string) ([]models.Account, error) {
	var accounts []models.Account
	err := conn.Select(&accounts, "SELECT id, balance, owner FROM accounts WHERE owner ILIKE $1", "%"+ownerName+"%")
	if err != nil {
		log.Printf("[ERROR] Search by owner '%s' failed: %v", ownerName, err)
		return nil, err
	}
	return accounts, err
}

func GetAccountByID(conn DBOrTx, id int) (models.Account, error) {
	var account models.Account
	err := conn.Get(&account, "SELECT id,balance,owner FROM accounts WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Account{}, errors.New("счет с таким ID не существует")
		}
		log.Printf("[ERROR] GetAccountByID(id=%d) failed: %v", id, err)
		return models.Account{}, err
	}

	return account, nil
}

func TopUpAccount(conn DBOrTx, accountID int, amount int64) error {
	_, err := conn.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, accountID)
	if err != nil {
		log.Printf("[ERROR] TopUP (id=%d, amount=%d) failed: %v", accountID, amount, err)
	}
	return err
}

func WithdrawAccount(conn DBOrTx, accountID int, amount int64) error {
	_, err := conn.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, accountID)
	if err != nil {
		log.Printf("[ERROR] Withdraw (id=%d, amount=%d) failed: %v", accountID, amount, err)
	}
	return err
}
