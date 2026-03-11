package repository

import (
	"github.com/MirzovalievShodmon/miniBank.git/internal/db"
	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
)

func GetAllAccounts() ([]models.Account, error) {
	var accounts []models.Account
	err := db.GetDBConnection().Select(&accounts, "SELECT id,balance,owner FROM accounts")
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func GetAccountByID(id int) (models.Account, error) {
	var account models.Account

	err := db.GetDBConnection().Get(&account, "SELECT id,balance,owner FROM accounts WHERE id = $1", id)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func TopUpAccount(accountID int, amount float64) error {
	_, err := db.GetDBConnection().Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, accountID)
	if err != nil {
		return err
	}

	return nil
}

func WithdrawAccount(accountID int, amount float64) error {
	_, err := db.GetDBConnection().Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, accountID)
	if err != nil {
		return err
	}

	return nil
}
