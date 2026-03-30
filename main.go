package main

import (
	"os"

	"github.com/MirzovalievShodmon/miniBank.git/internal/controller"
	"github.com/MirzovalievShodmon/miniBank.git/internal/db"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Str("module", "main").Msg("--- Запуск приложения miniBank ---")
	if err := db.InitConnection(); err != nil {
		log.Error().Str("module", "main").Err(err).Msg("Критическая Ошибка: База не подключена")
		return
	}
	log.Info().Str("module", "main").Msg("Соединение с базой данных установлено")

	if err := db.RunMigrations(); err != nil {
		log.Error().Str("module", "main").Err(err).Msg("Критическая Ошибка: Таблицы не созданы")
		return
	}
	log.Info().Str("module", "main").Msg("Миграции базы данных успешно применены")

	log.Info().Str("module", "main").Msg("Запуск веб-сервиса на порту :7556...")
	if err := controller.InitRoutes(); err != nil {
		log.Warn().
			Str("module", "main").
			Err(err).
			Msg("Предупреждение: Ошибка http-сервиса")
	}

	if err := db.CloseConnection(); err != nil {
		log.Error().Str("module", "main").Err(err).Msg("Ошибка при закрытии базы")
		return
	}
	log.Info().Msg("--- Приложение miniBank завершило работу ---")
}
