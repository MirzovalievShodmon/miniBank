package repository

import (
	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/rs/zerolog/log"
)

func CreateTransaction(conn DBOrTx, t models.Transaction) error {
	_, err := conn.Exec("INSERT INTO transactions (amount, type, account_id) VALUES ($1,$2,$3)", t.Amount, t.Type, t.AccountID)
	if err != nil {
		log.Error().
			Str("module", "repository").
			Str("table", "transactions").
			Int64("amount", t.Amount).
			Str("op_type", t.Type).
			Int("account_id", t.AccountID).
			Err(err).
			Msg("Не удалось записать транзакцию в БД")
	}
	return err
}

func GetTransactionsByAccountID(conn DBOrTx, accountID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := conn.Select(&transactions, "SELECT id, amount, type, created_at, account_id FROM transactions WHERE account_id = $1", accountID)
	if err != nil {
		log.Error().
			Str("module", "repository").
			Int("account_id", accountID).
			Err(err).
			Msg("Сбой при получении истории для конкретного аккаунта")
		return nil, err
	}
	return transactions, nil
}

func GetAllTransactions(conn DBOrTx) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := conn.Select(&transactions, "SELECT id,amount,type,created_at,account_id FROM transactions")
	if err != nil {
		log.Error().
			Str("module", "repository").
			Err(err).
			Msg("Сбой при получении всей истории транзакций")
		return nil, err
	}
	return transactions, nil
}
