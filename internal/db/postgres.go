package db

import (
	"time"

	"github.com/MirzovalievShodmon/miniBank.git/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var db *sqlx.DB

// InitConnection открытие подключения к БД с конфигурацией
func InitConnection(cfg *config.Config) error {
	connString := cfg.GetDatabaseConnectionString()

	log.Info().
		Str("module", "database").
		Str("host", cfg.Database.Host).
		Int("port", cfg.Database.Port).
		Str("database", cfg.Database.DBName).
		Msg("Попытка подключения к серверу PostgreSQL")

	dbConn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Error().
			Str("module", "database").
			Err(err).
			Msg("Сбой подключения к БД")
		return err
	}

	// Настройка пула соединений
	dbConn.SetMaxOpenConns(25)                 // Максимальное количество открытых соединений
	dbConn.SetMaxIdleConns(5)                  // Максимальное количество простаивающих соединений
	dbConn.SetConnMaxLifetime(5 * time.Minute) // Максимальное время жизни соединения

	// Проверка соединения
	if err := dbConn.Ping(); err != nil {
		log.Error().
			Str("module", "database").
			Err(err).
			Msg("Не удалось выполнить ping к БД")
		return err
	}

	db = dbConn
	log.Info().
		Str("module", "database").
		Str("db_name", cfg.Database.DBName).
		Msg("Успешное подключение к базе данных")
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
