package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Recovery простой middleware для обработки паник
// Панка - это критическая ошибка, которая останавливает выполнение программы
// Этот middleware ловит такие ошибки и не дает серверу упасть
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		// defer выполнится даже если произойдет паника
		defer func() {
			// recover() ловит панику и возвращает её значение
			// Если паники не было, recover() возвращает nil
			if err := recover(); err != nil {
				// Логируем ошибку для разработчика
				log.Error().
					Str("module", "recovery").
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path).
					Interface("panic", err).
					Msg("Поймана паника в обработчике")

				// Отправляем клиенту понятную ошибку вместо падения сервера
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Внутренняя ошибка сервера",
				})

				// Прерываем выполнение запроса
				c.Abort()
			}
		}()

		// Выполняем следующий обработчик
		c.Next()
	}
}
