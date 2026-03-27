package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes() error {
	r := gin.Default()

	r.GET("/ping", ping)

	accountsG := r.Group("/accounts", middleware)

	accountsG.GET("", getAllAccounts)

	accountsG.GET("/search", getAccountsByOwner)

	accountsG.POST("/:id/top-up", topUpAccount)

	accountsG.POST("/:id/withdraw", withdrawAccount)

	accountsG.POST("/transfer", transferAccount)

	accountsG.GET("/:id/transactions", getTransactionsByAccountID)

	r.GET("/transactions", getAllTransactions)

	return r.Run(":7556")
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}

func middleware(c *gin.Context) {
	fmt.Println("Я вызвался до")
	c.Next()
}
