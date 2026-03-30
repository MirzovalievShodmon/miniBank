package service

import (
	"github.com/MirzovalievShodmon/miniBank.git/internal/db"
	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/MirzovalievShodmon/miniBank.git/internal/repository"
	"github.com/rs/zerolog/log"
)

func GetAllTransactions() ([]models.Transaction, error) {
	conn := db.GetDBConnection()

	transactions, err := repository.GetAllTransactions(conn)

	if err != nil {
		log.Error().
			Str("module", "service").
			Err(err).
			Msg("Сбой при получении общего списка транзакций")
	}

	return transactions, err
}

func GetTransactionsByAccountID(accountID int) ([]models.Transaction, error) {
	conn := db.GetDBConnection()
	transactions, err := repository.GetTransactionsByAccountID(conn, accountID)

	if err != nil {
		log.Error().
			Str("module", "service").
			Int("account_id", accountID).
			Err(err).
			Msg("Сбой при получении истории счета")
	}
	return transactions, err
}
