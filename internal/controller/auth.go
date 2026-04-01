package controller

import (
	"net/http"

	"github.com/MirzovalievShodmon/miniBank.git/internal/config"
	"github.com/MirzovalievShodmon/miniBank.git/internal/middleware"
	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/MirzovalievShodmon/miniBank.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// registerUser обрабатывает регистрацию новых пользователей
func registerUser(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().
			Str("module", "controller").
			Err(err).
			Msg("Ошибка парсинга JSON при регистрации")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных",
		})
		return
	}

	log.Info().
		Str("module", "controller").
		Str("email", req.Email).
		Str("first_name", req.FirstName).
		Msg("Запрос на регистрацию")

	user, err := service.RegisterUser(req)
	if err != nil {
		log.Error().
			Str("module", "controller").
			Str("email", req.Email).
			Err(err).
			Msg("Ошибка регистрации пользователя")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Info().
		Str("module", "controller").
		Int("user_id", user.ID).
		Str("email", user.Email).
		Msg("Пользователь успешно зарегистрирован")

	c.JSON(http.StatusCreated, gin.H{
		"message": "Пользователь успешно зарегистрирован",
		"user":    user,
	})
}

// loginUser обрабатывает вход пользователей
func loginUser(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error().
				Str("module", "controller").
				Err(err).
				Msg("Ошибка парсинга JSON при входе")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Неверный формат данных",
			})
			return
		}

		log.Info().
			Str("module", "controller").
			Str("email", req.Email).
			Msg("Запрос на вход")

		authResponse, err := service.LoginUser(req, cfg)
		if err != nil {
			log.Error().
				Str("module", "controller").
				Str("email", req.Email).
				Err(err).
				Msg("Ошибка входа пользователя")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		log.Info().
			Str("module", "controller").
			Int("user_id", authResponse.User.ID).
			Str("email", authResponse.User.Email).
			Msg("Пользователь успешно авторизован")

		c.JSON(http.StatusOK, authResponse)
	}
}

// getProfile возвращает профиль текущего пользователя
func getProfile(c *gin.Context) {
	user, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Пользователь не авторизован",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// updateProfile обновляет профиль пользователя
func updateProfile(c *gin.Context) {
	user, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Пользователь не авторизован",
		})
		return
	}

	var updateReq struct {
		FirstName string `json:"first_name" binding:"omitempty,min=2"`
		LastName  string `json:"last_name" binding:"omitempty,min=2"`
	}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных",
		})
		return
	}

	// Обновляем только переданные поля
	if updateReq.FirstName != "" {
		user.FirstName = updateReq.FirstName
	}
	if updateReq.LastName != "" {
		user.LastName = updateReq.LastName
	}

	// TODO: Добавить вызов repository для обновления в БД

	log.Info().
		Str("module", "controller").
		Int("user_id", user.ID).
		Msg("Профиль пользователя обновлен")

	c.JSON(http.StatusOK, gin.H{
		"message": "Профиль обновлен",
		"user":    user,
	})
}
