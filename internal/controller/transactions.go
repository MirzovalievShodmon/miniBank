package controller

import (
	"net/http"

	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/MirzovalievShodmon/miniBank.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func getAllTransactions(c *gin.Context) {
	transactions, err := service.GetAllTransactions()
	if err != nil {
		log.Error().
			Str("module", "controller").
			Err(err).
			Msg("Сбой при получении всей истории транзакций")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(transactions) == 0 {
		c.JSON(200, []models.Transaction{})
		return
	}

	c.JSON(200, transactions)
}
