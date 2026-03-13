package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var db *sqlx.DB

func InitConnection(log zerolog.Logger) (err error) {
	// todo: поменял пароль на 1234, верни потом
	db, err = sqlx.Connect("postgres", "user=postgres password=1234 dbname=bank_db sslmode=disable")
	if err != nil {
		log.Error().Err(err).Msgf("error in initializing database connection")
		return err
	}

	return nil
}

func CloseConnection() error {
	err := db.Close()
	if err != nil {
		log.Error().Err(err).Msgf("error in initializing database connection")
		return err
	}

	return nil
}

func GetDBConnection() *sqlx.DB {
	return db
}
