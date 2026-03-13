package main

import (
	"github.com/MirzovalievShodmon/miniBank.git/internal/pkg/logger"

	"github.com/MirzovalievShodmon/miniBank.git/internal/controller"
	"github.com/MirzovalievShodmon/miniBank.git/internal/db"
)

// The Power of Justice
// Закончить правильную обработку ошибок
// gin-gonic
// swagger
// Добавить функцию перевод между абонентами
func main() {
	// Инициализация zerolog
	log := logger.InitLogger()
	log.Info().Msg("App started.")

	// Инициализация базы данных
	if err := db.InitConnection(log); err != nil {
		log.Error().Err(err).Msg("error during database connection initialization")
		return
	}
	defer func() {
		err := db.CloseConnection()
		if err != nil {
			log.Error().Err(err).Msg("error during database connection close")
			return
		}
	}()

	if err := db.RunMigrations(log); err != nil {
		log.Error().Err(err).Msg("Error during database migrations")
		return
	}

	controller.InitRoutes(log)
}
