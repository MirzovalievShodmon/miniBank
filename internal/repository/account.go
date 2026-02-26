package repository

import (
	"errors"

	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
)

func GetAllAccounts() ([]models.Account, error) {
	return accounts, nil
}

func GetAccountByID(id int) (models.Account, error) {
	for _, account := range accounts {
		if account.ID == id {
			return account, nil
		}
	}

	return models.Account{}, errors.New("нет счета с таким ID")
}

func TopUpAccount(accountID int, amount float64) error {
	for i := range accounts {
		if accounts[i].ID == accountID {
			accounts[i].Balance += amount
			return nil
		}
	}

	return errors.New("нет счета с таким ID")
}

func WithdrawAccount(accountID int, amount float64) error {
	for i := range accounts {
		if accounts[i].ID == accountID {
			accounts[i].Balance -= amount
			return nil
		}
	}

	return errors.New("нет счета с таким ID")
}
