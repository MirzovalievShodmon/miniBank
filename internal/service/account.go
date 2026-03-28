package service

import (
	"errors"
	"log"
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

func TopUpAccount(accountID int, amount int64) error {
	log.Printf("[SERVICE] Попытка пополнения счета ID %d на сумму %d", accountID, amount)

	conn := db.GetDBConnection()
	//Начинаем транзакцию
	tx, err := conn.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// проверка существует ли аккаунт
	if _, err = repository.GetAccountByID(tx, accountID); err != nil {
		return err
	}

	//пополнение аккаунта
	if err = repository.TopUpAccount(tx, accountID, amount); err != nil {
		return err
	}

	// сохранение записи о транзакции
	t := models.Transaction{
		AccountID: accountID,
		Amount:    amount,
		Type:      "Пополнение",
		CreatedAt: time.Now(),
	}
	if err = repository.CreateTransaction(tx, t); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	log.Printf("[SUCCES] Счет %d пополнен на %d", accountID, amount)
	return nil
}

func WithdrawAccount(accountID int, amount int64) error {
	log.Printf("[SERVICE] Попытка снятия со счета ID %d на сумму %d", accountID, amount)

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
		log.Printf("[WARN] Отказ в снятии: на счету %d недостаточно средств", accountID)
		return errors.New("недостаточно средств на счета")
	}

	//  снятие деньги c аккаунта
	if err = repository.WithdrawAccount(tx, accountID, amount); err != nil {
		return err
	}

	// сохранение записи о транзакции
	t := models.Transaction{
		AccountID: accountID,
		Amount:    amount,
		Type:      "Снятие",
		CreatedAt: time.Now(),
	}
	if err = repository.CreateTransaction(tx, t); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	log.Printf("[SUCCESS] Со счета %d снято %d", accountID, amount)
	return nil
}

func Transfer(fromID int, toID int, amount int64) (int64, error) {
	log.Printf("[SERVICE] Инициация перевода: от %d к %d на сумму %d", fromID, toID, amount)

	conn := db.GetDBConnection()
	tx, err := conn.Beginx() // начинаем транзакцию
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var fromAcc models.Account
	fromAcc, err = repository.GetAccountByID(tx, fromID)
	if err != nil {
		return 0, errors.New("отправитель не найден")
	}
	if fromAcc.Balance < amount {
		log.Printf("[WARN] Перевод отменен: у счета %d недостаточно денег", fromID)
		return 0, errors.New("недостаточно денег")
	}

	newBalance := fromAcc.Balance - amount

	if err := repository.WithdrawAccount(tx, fromID, amount); err != nil {
		return 0, err
	}

	if err := repository.TopUpAccount(tx, toID, amount); err != nil {
		return 0, err
	}

	// ПЕРЕДАЕМ tx! Теперь запись в лог под защитой транзакции
	t := models.Transaction{AccountID: fromID, Amount: amount, Type: "Перевод"}
	if err := repository.CreateTransaction(tx, t); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	log.Printf("[SUCCESS] Перевод выполнен: %d -> %d (%d)", fromID, toID, amount)
	return newBalance, nil
}
