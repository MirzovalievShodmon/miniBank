package db

import (
	"github.com/rs/zerolog/log"
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
	log.Info().Str("module", "database").Msg("Запуск миграций: создание таблицы accounts...")
	_, err := db.Exec(createTableAccountsDDL)
	if err != nil {
		log.Error().
			Str("module", "database").
			Str("table", "accounts").
			Err(err).
			Msg("Не удалось создать таблицу accounts")
		return err
	}
	log.Info().Str("module", "database").Str("table", "accounts").Msg("Таблица accounts готова")

	log.Info().Str("module", "database").Msg("Запуск миграций: создание таблицы transactions...")
	_, err = db.Exec(createTableTransactionsDDL)
	if err != nil {
		log.Error().
			Str("module", "database").
			Str("table", "transactions").
			Err(err).
			Msg("Не удалось создать таблицу transactions")
		return err
	}
	log.Info().Str("module", "database").Str("table", "transactions").Msg("Таблица transactions готова")

	return nil
}
