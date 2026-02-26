package repository

import "github.com/MirzovalievShodmon/miniBank.git/internal/models"

func GetAllTransactions() ([]models.Transaction, error) {
	return transacitons, nil
}

func CreateTransaction(t models.Transaction) error {
	t.ID = len(transacitons) + 1
	transacitons = append(transacitons, t)
	return nil
}
