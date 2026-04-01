package app

import (
	"os"

	"github.com/MirzovalievShodmon/miniBank.git/internal/config"
	"github.com/MirzovalievShodmon/miniBank.git/internal/controller"
	"github.com/MirzovalievShodmon/miniBank.git/internal/db"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Run Тут находится вся логика запуска приложения, которая была в main.go, чтобы main.go был максимально чистым и понятным
func Run() {
	// Инициализация логгера
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Str("module", "main").Msg("--- Запуск приложения miniBank ---")

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Str("module", "main").Err(err).Msg("Критическая ошибка: Не удалось загрузить конфигурацию")
		return
	}

	// Настройка уровня логирования
	switch cfg.App.LogLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Подключение к базе данных
	if err := db.InitConnection(cfg); err != nil {
		log.Fatal().Str("module", "main").Err(err).Msg("Критическая ошибка: База не подключена")
		return
	}
	log.Info().Str("module", "main").Msg("Соединение с базой данных установлено")

	// Выполнение миграций
	if err := db.RunMigrations(); err != nil {
		log.Fatal().Str("module", "main").Err(err).Msg("Критическая ошибка: Таблицы не созданы")
		return
	}
	log.Info().Str("module", "main").Msg("Миграции базы данных успешно применены")

	// Запуск веб-сервера
	log.Info().
		Str("module", "main").
		Str("port", cfg.Server.Port).
		Str("environment", cfg.App.Environment).
		Msg("Запуск веб-сервиса...")

	if err := controller.InitRoutes(cfg); err != nil {
		log.Error().
			Str("module", "main").
			Err(err).
			Msg("Критическая ошибка: HTTP-сервер не запустился")
	}

	// Закрытие соединения с базой данных
	if err := db.CloseConnection(); err != nil {
		log.Error().Str("module", "main").Err(err).Msg("Ошибка при закрытии базы")
		return
	}
	log.Info().Msg("--- Приложение miniBank завершило работу ---")
}
