package models

import "time"

type Transaction struct {
	ID        int       `db:"id" json:"id"`
	Amount    int64     `db:"amount" json:"amount"`
	Type      string    `db:"type" json:"type"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	AccountID int       `db:"account_id" json:"account_id"`
}
