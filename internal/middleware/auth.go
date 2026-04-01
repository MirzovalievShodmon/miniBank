package middleware

import (
	"net/http"
	"strings"

	"github.com/MirzovalievShodmon/miniBank.git/internal/config"
	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/MirzovalievShodmon/miniBank.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// AuthMiddleware middleware для проверки JWT токена
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Отсутствует токен авторизации",
			})
			c.Abort()
			return
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Неверный формат токена. Используйте: Bearer <token>",
			})
			c.Abort()
			return
		}

		token := parts[1]
		user, err := service.ValidateJWTToken(token, cfg)
		if err != nil {
			log.Warn().
				Str("module", "middleware").
				Str("client_ip", c.ClientIP()).
				Err(err).
				Msg("Невалидный токен в запросе")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Невалидный токен",
			})
			c.Abort()
			return
		}

		// Сохраняем информацию о пользователе в контексте
		c.Set("user", user)
		c.Set("user_id", user.ID)

		log.Info().
			Str("module", "middleware").
			Int("user_id", user.ID).
			Str("email", user.Email).
			Str("path", c.Request.URL.Path).
			Msg("Авторизованный запрос")

		c.Next()
	}
}

// GetCurrentUser возвращает текущего пользователя из контекста
func GetCurrentUser(c *gin.Context) (*models.User, bool) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		return nil, false
	}

	return user, true
}

// GetCurrentUserID возвращает ID текущего пользователя из контекста
func GetCurrentUserID(c *gin.Context) (int, bool) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		return 0, false
	}

	return userID, true
}

// RequireAuth требует авторизации для эндпойнта
func RequireAuth(cfg *config.Config) gin.HandlerFunc {
	return AuthMiddleware(cfg)
}
