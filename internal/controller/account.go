package controller

import (
	"net/http"
	"strconv"

	"github.com/MirzovalievShodmon/miniBank.git/internal/service"
	"github.com/gin-gonic/gin"
)

type patchBalance struct {
	Amount float64 `json:"amount"`
}

func getAllAccounts(c *gin.Context) {
	accounts, err := service.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func topUpAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный ID счета",
		})
		return
	}

	var req patchBalance

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат суммы",
		})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "сумма должна быть больше нуля",
		})
		return
	}

	err = service.TopUpAccount(id, req.Amount)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Баланс успешно пополнен",
	})
}

func WithdrawAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный ID счета",
		})
		return
	}

	var req patchBalance

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат суммы",
		})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "сумма должна быть больше нуля",
		})
		return
	}

	err = service.WithdrawAccount(id, req.Amount)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Деньги успешно сняты со счета",
	})
}
