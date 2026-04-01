package models

import "time"

type Account struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Balance   int64     `db:"balance" json:"balance"`
	Name      string    `db:"name" json:"name"` // Название счета (например "Основной счет")
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	IsActive  bool      `db:"is_active" json:"is_active"`
}

type CreateAccountRequest struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}
