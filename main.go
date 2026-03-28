package main

import (
	"log"

	"github.com/MirzovalievShodmon/miniBank.git/internal/controller"
	"github.com/MirzovalievShodmon/miniBank.git/internal/db"
)

func main() {
	log.Println("---  Запуск приложения miniBank ---")
	if err := db.InitConnection(); err != nil {
		log.Printf("Критическая Ошибка: База не пдключена: %v", err)
		return
	}
	log.Println("Соединение с базой данных установлено")

	if err := db.RunMigrations(); err != nil {
		log.Printf("Критическая Ошибка: Таблицы не созданы: %v", err)
		return
	}
	log.Println("Миграции базы данных успешно применены")

	log.Println("Запуск веб-сервиса на порту :7556...")
	if err := controller.InitRoutes(); err != nil {
		log.Printf("Предупреждение: Ошибка http-сервиса: %v", err)
	}

	if err := db.CloseConnection(); err != nil {
		log.Printf("Ошибка при закрытии базы: %v", err)
		return
	}
	log.Println("--- Приложение miniBank завершило работу ---")
}
