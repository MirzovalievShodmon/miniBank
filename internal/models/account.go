package models

type Account struct {
	ID      int     `db:"id"`
	Balance float64 `db:"balance"`
	Owner   string  `db:"owner"`
}
