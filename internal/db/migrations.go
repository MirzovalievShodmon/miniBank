package db

import (
	"github.com/rs/zerolog/log"
)

const (
	createTableUsersDDL = `CREATE TABLE IF NOT EXISTS users
(
	id SERIAL PRIMARY KEY,
	email VARCHAR(255) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	first_name VARCHAR(100) NOT NULL,
	last_name VARCHAR(100) NOT NULL,
	is_active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	createTableAccountsDDL = `CREATE TABLE IF NOT EXISTS accounts
(
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	balance BIGINT DEFAULT 0,
	name VARCHAR(100) NOT NULL,
	is_active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	createTableTransactionsDDL = `CREATE TABLE IF NOT EXISTS transactions
(
	id SERIAL PRIMARY KEY,
	amount BIGINT NOT NULL,
	type VARCHAR(50) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	account_id INT NOT NULL REFERENCES accounts (id) ON DELETE CASCADE
);`

	createIndexUserEmailDDL          = `CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`
	createIndexAccountUserDDL        = `CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);`
	createIndexTransactionAccountDDL = `CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);`
)

func RunMigrations() error {
	log.Info().Str("module", "database").Msg("Запуск миграций базы данных...")

	// Создание таблицы пользователей
	log.Info().Str("module", "database").Msg("Создание таблицы users...")
	_, err := db.Exec(createTableUsersDDL)
	if err != nil {
		log.Error().
			Str("module", "database").
			Str("table", "users").
			Err(err).
			Msg("Не удалось создать таблицу users")
		return err
	}
	log.Info().Str("module", "database").Str("table", "users").Msg("Таблица users готова")

	// Создание обновленной таблицы счетов
	log.Info().Str("module", "database").Msg("Создание таблицы accounts...")
	_, err = db.Exec(createTableAccountsDDL)
	if err != nil {
		log.Error().
			Str("module", "database").
			Str("table", "accounts").
			Err(err).
			Msg("Не удалось создать таблицу accounts")
		return err
	}
	log.Info().Str("module", "database").Str("table", "accounts").Msg("Таблица accounts готова")

	// Создание таблицы транзакций
	log.Info().Str("module", "database").Msg("Создание таблицы transactions...")
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

	// Создание индексов для производительности
	log.Info().Str("module", "database").Msg("Создание индексов...")

	_, err = db.Exec(createIndexUserEmailDDL)
	if err != nil {
		log.Warn().Str("module", "database").Err(err).Msg("Не удалось создать индекс для users.email")
	}

	_, err = db.Exec(createIndexAccountUserDDL)
	if err != nil {
		log.Warn().Str("module", "database").Err(err).Msg("Не удалось создать индекс для accounts.user_id")
	}

	_, err = db.Exec(createIndexTransactionAccountDDL)
	if err != nil {
		log.Warn().Str("module", "database").Err(err).Msg("Не удалось создать индекс для transactions.account_id")
	}

	log.Info().Str("module", "database").Msg("Миграции базы данных завершены успешно")
	return nil
}
