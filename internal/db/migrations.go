package db

import (
	"log"
)

const (
	createTableAccountsDDL = `CREATE TABLE IF NOT EXISTS accounts
(
	id SERIAL PRIMARY KEY,
	balance BIGINT DEFAULT 0,
	owner VARCHAR
);`

	createTableTransactionsDDL = `CREATE TABLE IF NOT EXISTS transactions
(
	id SERIAL PRIMARY KEY,
	amount BIGINT,
	type VARCHAR,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	account_id INT REFERENCES accounts (id)
);`
)

func RunMigrations() error {
	log.Println("[DB] Создание таблицы accounts...")
	_, err := db.Exec(createTableAccountsDDL)
	if err != nil {
		log.Printf("[ERROR] Не удалось создать таблицу accounts: %v", err)
		return err
	}
	log.Println("[DB] Таблтца accounts готова")

	log.Println("[DB] Создание таблицы transactions...")
	_, err = db.Exec(createTableTransactionsDDL)
	if err != nil {
		log.Printf("[ERROR] Не удалось создать таблицу transactions: %v", err)
		return err
	}
	log.Println("[DB] Таблица transactions готова")

	return nil
}
