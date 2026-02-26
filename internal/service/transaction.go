package service

import (
	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/MirzovalievShodmon/miniBank.git/internal/repository"
)

func GetAllTransactions() ([]models.Transaction, error) {
	return repository.GetAllTransactions()
}
