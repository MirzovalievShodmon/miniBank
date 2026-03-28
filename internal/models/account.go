package models

type Account struct {
	ID      int    `db:"id" json:"id"`
	Balance int64  `db:"balance" json:"balance"`
	Owner   string `db:"owner" json:"owner"`
}
