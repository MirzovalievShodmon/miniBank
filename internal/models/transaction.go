package models

import "time"

type Transaction struct {
	ID        int
	Amount    float64
	Type      string
	CreatedAt time.Time
	AccountID int
}
