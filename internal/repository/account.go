package repository

import (
	"database/sql"
	"errors"

	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/rs/zerolog/log"
)

type DBOrTx interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// CreateAccount создает новый счет в базе данных
func CreateAccount(conn DBOrTx, account models.Account) (models.Account, error) {
	query := `
		INSERT INTO accounts (user_id, name, balance, is_active, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
		RETURNING id, user_id, name, balance, is_active, created_at, updated_at`

	var newAccount models.Account
	err := conn.Get(&newAccount, query, account.UserID, account.Name, account.Balance, account.IsActive)
	if err != nil {
		log.Error().
			Str("module", "repository").
			Int("user_id", account.UserID).
			Err(err).
			Msg("Ошибка создания счета в БД")
		return models.Account{}, err
	}

	return newAccount, nil
}

// GetUserAccounts получает все счета пользователя
func GetUserAccounts(conn DBOrTx, userID int) ([]models.Account, error) {
	var accounts []models.Account
	query := `SELECT id, user_id, name, balance, is_active, created_at, updated_at FROM accounts WHERE user_id = $1 AND is_active = true`
	err := conn.Select(&accounts, query, userID)
	if err != nil {
		log.Error().
			Str("module", "repository").
			Int("user_id", userID).
			Err(err).
			Msg("Сбой при получении счетов пользователя")
		return nil, err
	}
	return accounts, nil
}

func GetAllAccounts(conn DBOrTx) ([]models.Account, error) {
	var accounts []models.Account
	query := `SELECT id, user_id, name, balance, is_active, created_at, updated_at FROM accounts WHERE is_active = true`
	err := conn.Select(&accounts, query)
	if err != nil {
		log.Error().
			Str("module", "repository").
			Err(err).
			Msg("Сбой при выполнении GetAllAccounts")
		return nil, err
	}
	return accounts, err
}

func GetAccountsByOwner(conn DBOrTx, ownerName string) ([]models.Account, error) {
	// Поиск по имени владельца через JOIN с таблицей users
	var accounts []models.Account
	query := `
		SELECT a.id, a.user_id, a.name, a.balance, a.is_active, a.created_at, a.updated_at 
		FROM accounts a 
		JOIN users u ON a.user_id = u.id 
		WHERE (u.first_name ILIKE $1 OR u.last_name ILIKE $1) AND a.is_active = true`

	err := conn.Select(&accounts, query, "%"+ownerName+"%")
	if err != nil {
		log.Error().
			Str("module", "repository").
			Str("search_term", ownerName).
			Err(err).
			Msg("Сбой при поиске аккаунтов по владельцу")
		return nil, err
	}
	return accounts, err
}

func GetAccountByID(conn DBOrTx, id int) (models.Account, error) {
	var account models.Account
	query := `SELECT id, user_id, name, balance, is_active, created_at, updated_at FROM accounts WHERE id = $1 AND is_active = true`
	err := conn.Get(&account, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Account{}, errors.New("счет с таким ID не существует")
		}
		log.Error().
			Str("module", "repository").
			Int("requested_id", id).
			Err(err).
			Msg("Критический сбой БД при поиске аккаунта")

		return models.Account{}, err
	}
	return account, nil
}

// GetAccountByIDAndUserID получает счет по ID и проверяет принадлежность пользователю
func GetAccountByIDAndUserID(conn DBOrTx, accountID, userID int) (models.Account, error) {
	var account models.Account
	query := `SELECT id, user_id, name, balance, is_active, created_at, updated_at FROM accounts WHERE id = $1 AND user_id = $2 AND is_active = true`
	err := conn.Get(&account, query, accountID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Account{}, errors.New("счет не найден или не принадлежит пользователю")
		}
		log.Error().
			Str("module", "repository").
			Int("account_id", accountID).
			Int("user_id", userID).
			Err(err).
			Msg("Ошибка при поиске счета по ID и user ID")

		return models.Account{}, err
	}
	return account, nil
}

func TopUpAccount(conn DBOrTx, accountID int, amount int64) error {
	_, err := conn.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, accountID)
	if err != nil {
		log.Error().
			Str("module", "repository").
			Int("id", accountID).
			Int64("amount", amount).
			Err(err).
			Msg("Сбой SQL при попытке ОБНОВЛЕНИЯ баланса (TopUp)")
	}
	return err
}

func WithdrawAccount(conn DBOrTx, accountID int, amount int64) error {
	_, err := conn.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, accountID)
	if err != nil {
		log.Error().
			Str("module", "repository").
			Int("id", accountID).
			Int64("amount", amount).
			Err(err).
			Msg("Сбой SQL при попытке ОБНОВЛЕНИЯ баланса (Withdraw)")
	}
	return err
}
