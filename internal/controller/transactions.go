package controller

import (
	"github.com/MirzovalievShodmon/miniBank.git/internal/service"
	"github.com/gin-gonic/gin"
)

func getAllTransactions(c *gin.Context) {
	transactions, err := service.GetAllTransactions()
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, transactions)
}
