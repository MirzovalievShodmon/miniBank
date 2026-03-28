package repository

import (
	"log"

	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
)

func CreateTransaction(conn DBOrTx, t models.Transaction) error {
	_, err := conn.Exec("INSERT INTO transactions (amount, type, account_id) VALUES ($1,$2,$3)", t.Amount, t.Type, t.AccountID)
	if err != nil {
		log.Printf("[ERROR] Не удалось записать транзакцию (Сумма: %d, Тип: %s, Аккаунт: %d): %v", t.Amount, t.Type, t.AccountID, err)
	}
	return err
}

func GetTransactionsByAccountID(conn DBOrTx, accountID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := conn.Select(&transactions, "SELECT id, amount, type, created_at, account_id FROM transactions WHERE account_id = $1", accountID)
	if err != nil {
		log.Printf("[ERROR] Сбой при получении истории для аккаунта %d: %v", accountID, err)
		return nil, err
	}
	return transactions, nil
}

func GetAllTransactions(conn DBOrTx) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := conn.Select(&transactions, "SELECT id,amount,type,created_at,account_id FROM transactions")
	if err != nil {
		log.Printf("[ERROR] Сбой при получении всей истории транзакций: %v", err)
		return nil, err
	}
	return transactions, nil
}
