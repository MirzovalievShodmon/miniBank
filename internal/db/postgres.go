package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var db *sqlx.DB

// открытие подключения к бд
func InitConnection() error {
	log.Info().Str("module", "database").Msg(" Попытка подключения к серверу PostgreSQL")
	dbConn, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=bank_db sslmode=disable")
	if err != nil {
		log.Error().
			Str("module", "database").
			Err(err).
			Msg("Сбой подключения к БД")
		return err
	}

	db = dbConn
	log.Info().
		Str("module", "database").
		Str("db_name", "bank_db").
		Msg("Успешное подключение к базе 'bank_db'")
	return nil
}

// закрытие подключения
func CloseConnection() error {
	log.Info().Str("module", "database").Msg("Закрытие соединение с БД")
	err := db.Close()
	if err != nil {
		log.Error().
			Str("module", "database").
			Err(err).Msg("Не удалось корректно закрыть БД")
		return err
	}
	log.Info().Str("module", "database").Msg("Соединение с базой успешно закрыто")
	return nil
}

// получение объект соединения
func GetDBConnection() *sqlx.DB {
	return db
}
