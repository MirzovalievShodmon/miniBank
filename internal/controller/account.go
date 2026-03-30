package controller

import (
	"net/http"
	"strconv"

	"github.com/MirzovalievShodmon/miniBank.git/internal/models"
	"github.com/MirzovalievShodmon/miniBank.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type patchBalance struct {
	Amount int64 `json:"amount"`
}

type transferRequest struct {
	FromID int   `json:"from_id"`
	ToID   int   `json:"to_id"`
	Amount int64 `json:"amount"`
}

func getAllAccounts(c *gin.Context) {
	accounts, err := service.GetAllAccounts()
	if err != nil {
		log.Error().
			Str("module", "controller").
			Err(err).
			Msg("Сбой при получении списка всех счетов")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(accounts) == 0 {
		c.JSON(200, []models.Account{})
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
		log.Error().
			Str("module", "controller").
			Str("search_term", name).
			Err(err).
			Msg("Сбой при поиске аккаунтов по владельцу")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка поиска",
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
		log.Error().
			Str("module", "controller").
			Int("account_id", id).
			Err(err).
			Msg("Ошибка парсинга JSON при пополнении")
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

	log.Info().
		Str("module", "controller").
		Int("account_id", id).
		Int64("сумма", req.Amount).
		Msg("Запрос на пополнение принят")
	err = service.TopUpAccount(id, req.Amount)
	if err != nil {
		log.Error().
			Str("module", "controller").
			Int("account_id", id).
			Err(err).
			Msg("Сбой операции пополнения в сервисе")
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
		log.Error().
			Str("module", "controller").
			Int("account_id", id).
			Err(err).
			Msg("Ошибка парсинга JSON при снятии")
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

	log.Info().
		Str("module", "controller").
		Int("account_id", id).
		Int64("сумма", req.Amount).
		Msg("Запрос на снятие принят")
	err = service.WithdrawAccount(id, req.Amount)
	if err != nil {
		log.Error().
			Str("module", "controller").
			Int("account_id", id).
			Err(err).
			Msg("Сбой операции снятия в сервисе")
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
		log.Error().
			Str("module", "controller").
			Err(err).
			Msg("Ошибка парсинга JSON при переводе")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат JSON",
		})
		return
	}

	log.Info().
		Str("module", "controller").
		Int("from_id", req.FromID).
		Int("to_id", req.ToID).
		Int64("сумма", req.Amount).
		Msg("Запрос на перевод принят")
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

	newBalance, err := service.Transfer(req.FromID, req.ToID, req.Amount)
	if err != nil {
		log.Error().
			Str("module", "controller").
			Int("from_id", req.FromID).
			Int("to_id", req.ToID).
			Err(err).
			Msg("Сбой перевода в сервисе")
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "Перевод выполнен успешно",
		"sender_new_balance": newBalance,
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
		log.Error().
			Str("module", "controller").
			Int("account_id", id). // Логируем ID, чтобы знать, где проблема
			Err(err).
			Msg("Сбой при получении истории транзакций")
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
