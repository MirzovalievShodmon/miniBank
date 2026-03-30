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

func GetAllAccounts(conn DBOrTx) ([]models.Account, error) {
	var accounts []models.Account
	err := conn.Select(&accounts, "SELECT id,balance,owner FROM accounts")
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
	var accounts []models.Account
	err := conn.Select(&accounts, "SELECT id, balance, owner FROM accounts WHERE owner ILIKE $1", "%"+ownerName+"%")
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
	err := conn.Get(&account, "SELECT id,balance,owner FROM accounts WHERE id = $1", id)
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
