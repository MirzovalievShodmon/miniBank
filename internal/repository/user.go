package repository

import (
	"database/sql"
	"errors"

	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/rs/zerolog/log"
)

// CreateUser создает нового пользователя в базе данных
func CreateUser(conn DBOrTx, user models.User) (models.User, error) {
	query := `
		INSERT INTO users (email, password, first_name, last_name, is_active, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
		RETURNING id, email, first_name, last_name, is_active, created_at, updated_at`

	var newUser models.User
	err := conn.Get(&newUser, query, user.Email, user.Password, user.FirstName, user.LastName, user.IsActive)
	if err != nil {
		log.Error().
			Str("module", "repository").
			Str("email", user.Email).
			Err(err).
			Msg("Ошибка создания пользователя в БД")
		return models.User{}, err
	}

	return newUser, nil
}

// GetUserByEmail находит пользователя по email
func GetUserByEmail(conn DBOrTx, email string) (models.User, error) {
	var user models.User
	query := `SELECT id, email, password, first_name, last_name, is_active, created_at, updated_at FROM users WHERE email = $1`

	err := conn.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("пользователь с таким email не найден")
		}
		log.Error().
			Str("module", "repository").
			Str("email", email).
			Err(err).
			Msg("Ошибка поиска пользователя по email")
		return models.User{}, err
	}

	return user, nil
}

// GetUserByID находит пользователя по ID
func GetUserByID(conn DBOrTx, userID int) (models.User, error) {
	var user models.User
	query := `SELECT id, email, password, first_name, last_name, is_active, created_at, updated_at FROM users WHERE id = $1`

	err := conn.Get(&user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("пользователь не найден")
		}
		log.Error().
			Str("module", "repository").
			Int("user_id", userID).
			Err(err).
			Msg("Ошибка поиска пользователя по ID")
		return models.User{}, err
	}

	return user, nil
}

// UpdateUser обновляет информацию о пользователе
func UpdateUser(conn DBOrTx, user models.User) error {
	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, is_active = $3, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $4`

	_, err := conn.Exec(query, user.FirstName, user.LastName, user.IsActive, user.ID)
	if err != nil {
		log.Error().
			Str("module", "repository").
			Int("user_id", user.ID).
			Err(err).
			Msg("Ошибка обновления пользователя")
		return err
	}

	return nil
}

// ChangeUserPassword изменяет пароль пользователя
func ChangeUserPassword(conn DBOrTx, userID int, hashedPassword string) error {
	query := `UPDATE users SET password = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	_, err := conn.Exec(query, hashedPassword, userID)
	if err != nil {
		log.Error().
			Str("module", "repository").
			Int("user_id", userID).
			Err(err).
			Msg("Ошибка изменения пароля пользователя")
		return err
	}

	return nil
}
