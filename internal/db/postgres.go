package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func InitConnection() error {
	dbConn, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=bank_db sslmode=disable")
	if err != nil {
		return err
	}

	db = dbConn
	return nil
}

func CloseConnection() error {
	err := db.Close()
	if err != nil {
		return err
	}

	return nil
}

func GetDBConnection() *sqlx.DB {
	return db
}
