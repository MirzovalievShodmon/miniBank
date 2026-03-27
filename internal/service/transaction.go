package service

import (
	"github.com/MirzovalievShodmon/miniBank.git/internal/db"
	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/MirzovalievShodmon/miniBank.git/internal/repository"
)

func GetAllTransactions() ([]models.Transaction, error) {
	conn := db.GetDBConnection()
	return repository.GetAllTransactions(conn)
}

func GetTransactionsByAccountID(accountID int) ([]models.Transaction, error) {
	conn := db.GetDBConnection()
	return repository.GetTransactionsByAccountID(conn, accountID)
}
