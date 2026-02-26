package repository

import "github.com/MirzovalievShodmon/miniBank.git/internal/models"

var (
	accounts = []models.Account{
		{
			ID:      1,
			Balance: 1000,
			Owner:   "Mirzovaliev Shodmon",
		},
		{
			ID:      2,
			Balance: 2000,
			Owner:   "Rahmonov Saidvaly",
		},
		{
			ID:      3,
			Balance: 1500,
			Owner:   "Dovudova Maryam",
		},
	}

	transacitons = []models.Transaction{}
)
