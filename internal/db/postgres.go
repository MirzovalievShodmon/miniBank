package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

// открытие подключения к бд
func InitConnection() error {
	log.Println("[DB] Попытка подключения к серверу PostgreSQL...")
	dbConn, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=bank_db sslmode=disable")
	if err != nil {
		log.Printf("[ERROR] Сбой подключения к БД: %v", err)
		return err
	}

	db = dbConn
	log.Println("[DB] Успешное подключение к базе 'bank_db'")
	return nil
}

// закрытие подключения
func CloseConnection() error {
	log.Println("[DB] Закрытие соединение с базой данных...")
	err := db.Close()
	if err != nil {
		log.Printf("[ERROR] Не удалось корректно закрыть БД: %v", err)
		return err
	}
	log.Println("[DB] Соединение с базой успешно закрыто")
	return nil
}

// получение объект соединения
func GetDBConnection() *sqlx.DB {
	return db
}
