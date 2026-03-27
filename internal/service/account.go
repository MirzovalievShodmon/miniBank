package service

import (
	"errors"
	"time"

	"github.com/MirzovalievShodmon/miniBank.git/internal/db"
	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/MirzovalievShodmon/miniBank.git/internal/repository"
)

func GetAllAccounts() ([]models.Account, error) {
	conn := db.GetDBConnection()
	return repository.GetAllAccounts(conn)
}

func GetAccountsByOwner(ownerName string) ([]models.Account, error) {
	conn := db.GetDBConnection()
	return repository.GetAccountsByOwner(conn, ownerName)
}

func TopUpAccount(accountID int, amount float64) error {
	conn := db.GetDBConnection()
	//Начинаем транзакцию
	tx, err := conn.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// проверка существует ли аккаунт
	_, err = repository.GetAccountByID(tx, accountID)
	if err != nil {
		return err
	}

	//пополнение аккаунта
	err = repository.TopUpAccount(tx, accountID, amount)
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
	err = repository.CreateTransaction(tx, t)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func WithdrawAccount(accountID int, amount float64) error {
	conn := db.GetDBConnection()
	// Начинаем транзакцию
	tx, err := conn.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// проверка существует ли аккаунт
	account, err := repository.GetAccountByID(tx, accountID)
	if err != nil {
		return err
	}

	if account.Balance < amount {
		return errors.New("недостаточно средств на счета")
	}

	//  снятие деньги c аккаунта
	err = repository.WithdrawAccount(tx, accountID, amount)
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
	err = repository.CreateTransaction(tx, t)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func Transfer(fromID int, toID int, amount float64) error {
	conn := db.GetDBConnection()
	tx, err := conn.Beginx() // начинаем транзакцию
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var fromAcc models.Account
	fromAcc, err = repository.GetAccountByID(tx, fromID)
	if err != nil {
		return errors.New("отправитель не найден")
	}
	if fromAcc.Balance < amount {
		return errors.New("недостаточно денег")
	}

	if err := repository.WithdrawAccount(tx, fromID, amount); err != nil {
		return err
	}

	if err := repository.TopUpAccount(tx, toID, amount); err != nil {
		return err
	}

	// ПЕРЕДАЕМ tx! Теперь запись в лог под защитой транзакции
	t := models.Transaction{AccountID: fromID, Amount: amount, Type: "Перевод"}
	if err := repository.CreateTransaction(tx, t); err != nil {
		return err
	}

	return tx.Commit()
}
