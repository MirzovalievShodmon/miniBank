package controller

import (
	"net/http"

	"github.com/MirzovalievShodmon/miniBank.git/internal/config"
	"github.com/MirzovalievShodmon/miniBank.git/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func InitRoutes(cfg *config.Config) error {
	// Настройка режима Gin в зависимости от окружения
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	// Добавляем middleware для логирования запросов и восстановления после ошибок
	r.Use(middleware.RequestLogging())
	r.Use(middleware.Recovery())

	// Публичные endpoints (без авторизации)
	r.GET("/ping", ping)
	r.GET("/health", healthCheck)

	// Группа для аутентификации
	authG := r.Group("/auth")
	authG.POST("/register", registerUser)
	authG.POST("/login", loginUser(cfg))

	// Защищенная группа для профиля пользователя
	profileG := r.Group("/profile", middleware.RequireAuth(cfg))
	profileG.GET("", getProfile)
	profileG.PUT("", updateProfile)

	// Защищенная группа для счетов (требует авторизации)
	accountsG := r.Group("/accounts", middleware.RequireAuth(cfg))
	accountsG.GET("", getAllAccounts)
	accountsG.POST("", createAccount)
	accountsG.GET("/search", getAccountsByOwner)
	accountsG.POST("/:id/top-up", topUpAccount)
	accountsG.POST("/:id/withdraw", withdrawAccount)
	accountsG.POST("/transfer", transferAccount)
	accountsG.GET("/:id/transactions", getTransactionsByAccountID)

	// Защищенная группа для транзакций (требует авторизации)
	transactionsG := r.Group("/transactions", middleware.RequireAuth(cfg))
	transactionsG.GET("", getAllTransactions)

	// Запуск сервера
	serverAddress := ":" + cfg.Server.Port
	log.Info().
		Str("module", "controller").
		Str("address", serverAddress).
		Str("environment", cfg.App.Environment).
		Msg("Запуск HTTP сервера")

	return r.Run(serverAddress)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "miniBank",
		"version": "1.0.0",
	})
}
