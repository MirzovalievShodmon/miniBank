package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	App      AppConfig
	Auth     AuthConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  int // в секундах
	WriteTimeout int // в секундах
	IdleTimeout  int // в секундах
}

type AppConfig struct {
	Environment string // development, staging, production
	LogLevel    string // debug, info, warn, error
}

type AuthConfig struct {
	JWTSecret        string // секретный ключ для JWT токенов
	TokenExpireHours int    // время жизни токена в часах
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() (*Config, error) {
	log.Info().Str("module", "config").Msg("Загрузка конфигурации из переменных окружения")

	// Попытка загрузить .env файл (игнорируем ошибку, если файл не найден)
	loadEnvFile()

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnvWithDefault("DB_HOST", "localhost"),
			Port:     getEnvAsIntWithDefault("DB_PORT", 5432),
			User:     getEnvWithDefault("DB_USER", "postgres"),
			Password: getEnvWithDefault("DB_PASSWORD", "postgres"),
			DBName:   getEnvWithDefault("DB_NAME", "bank_db"),
			SSLMode:  getEnvWithDefault("DB_SSL_MODE", "disable"),
		},
		Server: ServerConfig{
			Port:         getEnvWithDefault("SERVER_PORT", "7556"),
			ReadTimeout:  getEnvAsIntWithDefault("SERVER_READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsIntWithDefault("SERVER_WRITE_TIMEOUT", 10),
			IdleTimeout:  getEnvAsIntWithDefault("SERVER_IDLE_TIMEOUT", 120),
		},
		App: AppConfig{
			Environment: getEnvWithDefault("APP_ENV", "development"),
			LogLevel:    getEnvWithDefault("LOG_LEVEL", "info"),
		},
		Auth: AuthConfig{
			JWTSecret:        getEnvWithDefault("JWT_SECRET", "your-256-bit-secret-key-change-this-in-production!"),
			TokenExpireHours: getEnvAsIntWithDefault("JWT_EXPIRE_HOURS", 24),
		},
	}

	// Валидация критических параметров
	if config.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD не может быть пустым")
	}

	if config.Database.User == "" {
		return nil, fmt.Errorf("DB_USER не может быть пустым")
	}

	log.Info().
		Str("module", "config").
		Str("db_host", config.Database.Host).
		Int("db_port", config.Database.Port).
		Str("db_name", config.Database.DBName).
		Str("server_port", config.Server.Port).
		Str("environment", config.App.Environment).
		Msg("Конфигурация успешно загружена")

	return config, nil
}

// GetDatabaseConnectionString возвращает строку подключения к БД
func (c *Config) GetDatabaseConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// getEnvWithDefault возвращает значение переменной окружения или дефолтное значение
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsIntWithDefault возвращает значение переменной окружения как int или дефолтное значение
func getEnvAsIntWithDefault(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Warn().
			Str("key", key).
			Str("value", valueStr).
			Int("default", defaultValue).
			Msg("Не удается распарсить переменную окружения как int, используется значение по умолчанию")
		return defaultValue
	}

	return value
}

// loadEnvFile загружает переменные из .env файла (простая реализация без внешних зависимостей)
func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		// .env файл не найден - это нормально
		return
	}
	defer file.Close()

	log.Info().Str("module", "config").Msg(".env файл найден и загружается")

	buf := make([]byte, 1024)
	content := ""

	for {
		n, err := file.Read(buf)
		if n == 0 || err != nil {
			break
		}
		content += string(buf[:n])
	}

	// Разбираем строки
	for _, line := range splitLines(content) {
		line = trimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// Ищем знак равенства
		eqIndex := -1
		for i, char := range line {
			if char == '=' {
				eqIndex = i
				break
			}
		}

		if eqIndex == -1 {
			continue
		}

		key := trimSpace(line[:eqIndex])
		value := trimSpace(line[eqIndex+1:])

		if len(key) > 0 {
			os.Setenv(key, value)
		}
	}
}

// splitLines разбивает строку на строки
func splitLines(s string) []string {
	var lines []string
	var currentLine string

	for _, char := range s {
		if char == '\n' || char == '\r' {
			if len(currentLine) > 0 {
				lines = append(lines, currentLine)
				currentLine = ""
			}
		} else {
			currentLine += string(char)
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

// trimSpace удаляет пробелы в начале и конце строки
func trimSpace(s string) string {
	start := 0
	end := len(s)

	// Убираем пробелы в начале
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}

	// Убираем пробелы в конце
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}

	return s[start:end]
}
