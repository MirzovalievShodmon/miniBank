package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/MirzovalievShodmon/miniBank.git/internal/config"
	"github.com/MirzovalievShodmon/miniBank.git/internal/db"
	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/MirzovalievShodmon/miniBank.git/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength = 6
	TokenExpireHours  = 24
)

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// RegisterUser регистрирует нового пользователя
func RegisterUser(req models.RegisterRequest) (*models.User, error) {
	log.Info().
		Str("module", "service").
		Str("email", req.Email).
		Msg("Начало регистрации пользователя")

	// Проверяем, что пользователь с таким email не существует
	conn := db.GetDBConnection()
	existingUser, _ := repository.GetUserByEmail(conn, req.Email)
	if existingUser.ID != 0 {
		log.Warn().
			Str("module", "service").
			Str("email", req.Email).
			Msg("Попытка регистрации с существующим email")
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().
			Str("module", "service").
			Err(err).
			Msg("Ошибка хэширования пароля")
		return nil, errors.New("ошибка обработки пароля")
	}

	// Создаем пользователя
	user := models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	// Сохраняем в базу
	newUser, err := repository.CreateUser(conn, user)
	if err != nil {
		log.Error().
			Str("module", "service").
			Str("email", req.Email).
			Err(err).
			Msg("Ошибка создания пользователя в БД")
		return nil, err
	}

	log.Info().
		Str("module", "service").
		Int("user_id", newUser.ID).
		Str("email", newUser.Email).
		Msg("Пользователь успешно зарегистрирован")

	return &newUser, nil
}

// LoginUser аутентифицирует пользователя и возвращает токен
func LoginUser(req models.LoginRequest, cfg *config.Config) (*models.AuthResponse, error) {
	log.Info().
		Str("module", "service").
		Str("email", req.Email).
		Msg("Попытка входа пользователя")

	conn := db.GetDBConnection()
	user, err := repository.GetUserByEmail(conn, req.Email)
	if err != nil {
		log.Warn().
			Str("module", "service").
			Str("email", req.Email).
			Msg("Пользователь не найден")
		return nil, errors.New("неверный email или пароль")
	}

	// Проверяем активность пользователя
	if !user.IsActive {
		log.Warn().
			Str("module", "service").
			Int("user_id", user.ID).
			Msg("Попытка входа заблокированного пользователя")
		return nil, errors.New("аккаунт заблокирован")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Warn().
			Str("module", "service").
			Str("email", req.Email).
			Msg("Неверный пароль")
		return nil, errors.New("неверный email или пароль")
	}

	// Генерируем JWT токен
	token, err := generateJWTToken(user, cfg)
	if err != nil {
		log.Error().
			Str("module", "service").
			Int("user_id", user.ID).
			Err(err).
			Msg("Ошибка генерации JWT токена")
		return nil, errors.New("ошибка генерации токена")
	}

	log.Info().
		Str("module", "service").
		Int("user_id", user.ID).
		Str("email", user.Email).
		Msg("Пользователь успешно авторизован")

	// Убираем пароль из ответа
	user.Password = ""

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// generateJWTToken генерирует JWT токен для пользователя
func generateJWTToken(user models.User, cfg *config.Config) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(cfg.Auth.TokenExpireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "miniBank",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(cfg.Auth.JWTSecret))
}

// getJWTSecret возвращает секрет для JWT токенов
func getJWTSecret() string {
	// В production это должен быть настоящий секрет из переменных окружения
	// Для разработки используем фиксированный секрет (НЕ РЕКОМЕНДУЕТСЯ для production!)
	return "your-256-bit-secret-key-change-this-in-production!"
}

// ValidateJWTToken проверяет JWT токен и возвращает пользователя
func ValidateJWTToken(tokenString string, cfg *config.Config) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(cfg.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Получаем свежую информацию о пользователе из БД
		conn := db.GetDBConnection()
		user, err := repository.GetUserByID(conn, claims.UserID)
		if err != nil {
			return nil, errors.New("пользователь не найден")
		}

		if !user.IsActive {
			return nil, errors.New("пользователь заблокирован")
		}

		// Убираем пароль из ответа
		user.Password = ""
		return &user, nil
	}

	return nil, errors.New("невалидный токен")
}
