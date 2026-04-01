package service

import (
	"errors"
	"time"

	"github.com/MirzovalievShodmon/miniBank.git/internal/db"
	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/MirzovalievShodmon/miniBank.git/internal/repository"
	"github.com/rs/zerolog/log"
)

// CreateAccount создает новый счет для пользователя
func CreateAccount(userID int, accountName string) (*models.Account, error) {
	log.Info().
		Str("module", "service").
		Int("user_id", userID).
		Str("account_name", accountName).
		Msg("Создание нового счета")

	conn := db.GetDBConnection()

	account := models.Account{
		UserID:   userID,
		Name:     accountName,
		Balance:  0,
		IsActive: true,
	}

	newAccount, err := repository.CreateAccount(conn, account)
	if err != nil {
		log.Error().
			Str("module", "service").
			Int("user_id", userID).
			Err(err).
			Msg("Ошибка создания счета")
		return nil, err
	}

	log.Info().
		Str("module", "service").
		Int("user_id", userID).
		Int("account_id", newAccount.ID).
		Msg("Счет успешно создан")

	return &newAccount, nil
}

func GetAllAccounts() ([]models.Account, error) {
	conn := db.GetDBConnection()
	accounts, err := repository.GetAllAccounts(conn)
	if err != nil {
		log.Error().Str("module", "service").Err(err).Msg("Сбой получения списка счетов")
	}
	return accounts, err
}

func GetAccountsByOwner(ownerName string) ([]models.Account, error) {
	conn := db.GetDBConnection()
	accounts, err := repository.GetAccountsByOwner(conn, ownerName)
	if err != nil {
		log.Error().Str("module", "service").Str("owner", ownerName).Err(err).Msg("Сбой поиска по владельцу")
	}
	return accounts, err
}

func TopUpAccount(accountID int, amount int64) error {
	log.Info().
		Str("module", "service").
		Int("account_id", accountID).
		Int64("amount", amount).
		Msg("Начало пополнения счета")
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

	log.Info().Str("module", "service").Int("account_id", accountID).Msg("Пополнение успешно завершено")
	return nil
}

func WithdrawAccount(accountID int, amount int64) error {
	log.Info().
		Str("module", "service").
		Int("account_id", accountID).
		Int64("amount", amount).
		Msg("Начало снятия средств")

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
		log.Warn().
			Str("module", "service").
			Int("account_id", accountID).
			Int64("balance", account.Balance).
			Int64("requested", amount).
			Msg("Отказ: недостаточно средств")
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
		log.Error().Str("module", "service").Int("account_id", accountID).Err(err).Msg("Ошибка при Commit пополнения")
		return err
	}

	log.Info().Str("module", "service").Int("account_id", accountID).Msg("Снятие успешно завершено")
	return nil
}

func Transfer(fromID int, toID int, amount int64) (int64, error) {
	log.Info().
		Str("module", "service").
		Int("from_id", fromID).
		Int("to_id", toID).
		Int64("amount", amount).
		Msg("Инициация перевода")

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
		log.Warn().
			Str("module", "service").
			Int("from_id", fromID).
			Msg("Перевод отменен: недостаточно денег")
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

	log.Info().
		Str("module", "service").
		Int("from_id", fromID).
		Int("to_id", toID).
		Msg("Перевод успешно выполнен")

	return newBalance, nil
}
