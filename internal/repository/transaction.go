package repository

import (
	"github.com/MirzovalievShodmon/miniBank.git/internal/db"
	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
)

func CreateTransaction(t models.Transaction) error {
	_, err := db.GetDBConnection().Exec("INSERT INTO transactions (amount, type,account_id) VALUES ($1,$2,$3)", t.Amount, t.Type, t.AccountID)
	if err != nil {
		return err
	}

	return nil
}

func GetAllTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := db.GetDBConnection().Select(&transactions, "SELECT id,amount,type,created_at,account_id FROM transactions")
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
