package service

import (
	"errors"
	"time"

	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/MirzovalievShodmon/miniBank.git/internal/repository"
)

func GetAllAccounts() ([]models.Account, error) {
	return repository.GetAllAccounts()
}

func TopUpAccount(accountID int, amount float64) error {
	// проверка существует ли аккаунт
	_, err := repository.GetAccountByID(accountID)
	if err != nil {
		return err
	}

	//пополнение аккаунта
	err = repository.TopUpAccount(accountID, amount)
	if err != nil {
		return err
	}

	// сохранение записи о транзакции
	t := models.Transaction{
		AccountID: accountID,
		Amount:    amount,
		Type:      "Пополнение",
		CreatedAt: time.Now(),
	}
	err = repository.CreateTransaction(t)
	if err != nil {
		return err
	}

	return nil
}

func WithdrawAccount(accountID int, amount float64) error {
	// проверка существует ли аккаунт
	account, err := repository.GetAccountByID(accountID)
	if err != nil {
		return err
	}

	if account.Balance < amount {
		return errors.New("недостаточно средст на счета")
	}

	// пополнение аккаунта
	err = repository.WithdrawAccount(accountID, amount)
	if err != nil {
		return err
	}

	// сохранение записи о транзакции
	t := models.Transaction{
		AccountID: accountID,
		Amount:    amount,
		Type:      "Снятие",
		CreatedAt: time.Now(),
	}
	err = repository.CreateTransaction(t)
	if err != nil {
		return err
	}
	
	return nil
}
