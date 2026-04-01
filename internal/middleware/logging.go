package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RequestLogging простой middleware для логирования HTTP запросов
func RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Логируем входящий запрос
		log.Info().
			Str("module", "middleware").
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("client_ip", c.ClientIP()).
			Msg("Входящий HTTP запрос")

		// Выполняем обработчик
		c.Next()

		// Логируем ответ
		log.Info().
			Str("module", "middleware").
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Msg("HTTP ответ отправлен")
	}
}
