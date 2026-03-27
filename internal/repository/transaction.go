package repository

import (
	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
)

func CreateTransaction(conn DBOrTx, t models.Transaction) error {
	_, err := conn.Exec("INSERT INTO transactions (amount, type, account_id) VALUES ($1,$2,$3)", t.Amount, t.Type, t.AccountID)
	return err
}

func GetTransactionsByAccountID(conn DBOrTx, accountID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := conn.Select(&transactions, "SELECT id, amount, type, created_at, account_id FROM transactions WHERE account_id = $1", accountID)
	return transactions, err
}

func GetAllTransactions(conn DBOrTx) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := conn.Select(&transactions, "SELECT id,amount,type,created_at,account_id FROM transactions")
	return transactions, err
}
