package models

import "time"

type Transaction struct {
	ID        int       `db:"id"`
	Amount    float64   `db:"amount"`
	Type      string    `db:"type"`
	CreatedAt time.Time `db:"created_at"`
	AccountID int       `db:"account_id"`
}
