package controller

import (
	"net/http"
	"strconv"

	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/MirzovalievShodmon/miniBank.git/internal/service"
	"github.com/gin-gonic/gin"
)

type patchBalance struct {
	Amount float64 `json:"amount"`
}

type transferRequest struct {
	FromID int     `json:"from_id"`
	ToID   int     `json:"to_id"`
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

func getAccountsByOwner(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "параметр name не может быть пустым",
		})
		return
	}

	accounts, err := service.GetAccountsByOwner(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка поска",
		})
		return
	}

	if len(accounts) == 0 {
		c.JSON(http.StatusOK, []models.Account{})
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

func withdrawAccount(c *gin.Context) {
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

func transferAccount(c *gin.Context) {
	var req transferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат JSON",
		})
		return
	}

	if req.FromID == req.ToID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "нельзя перевести деньги самому себе",
		})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "сумма должна быть больше нуля",
		})
		return
	}

	if err := service.Transfer(req.FromID, req.ToID, req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Перевод выполнен успешно",
	})
}

func getTransactionsByAccountID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный ID счета",
		})
		return
	}
	history, err := service.GetTransactionsByAccountID(id)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "ошибка при получении истории",
		})
		return
	}
	if len(history) == 0 {
		c.JSON(200, []models.Transaction{})
		return
	}
	c.JSON(200, history)
}
