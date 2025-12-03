package account

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	useCase UseCase
	logger  *zap.Logger
}

func NewHandler(useCase UseCase, logger *zap.Logger) *Handler {
	return &Handler{useCase: useCase, logger: logger}
}

// Create godoc
// @Summary      Создать аккаунт
// @Description  Создаёт новый торговый аккаунт для пользователя
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        request body CreateAccountRequest true "Данные аккаунта"
// @Success      201 {object} Account
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /accounts [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.useCase.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetByID godoc
// @Summary      Получить аккаунт по ID
// @Description  Возвращает аккаунт по его идентификатору
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id path int true "ID аккаунта"
// @Success      200 {object} Account
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /accounts/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	account, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// List godoc
// @Summary      Список аккаунтов
// @Description  Возвращает список аккаунтов с пагинацией и фильтрами
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        user_id query int false "Фильтр по ID пользователя"
// @Param        account_type query string false "Фильтр по типу аккаунта (master/investor)"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество записей на странице" default(20)
// @Success      200 {object} AccountListResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /accounts [get]
func (h *Handler) List(c *gin.Context) {
	var filter AccountFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.useCase.List(c.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to list accounts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Update godoc
// @Summary      Обновить аккаунт
// @Description  Обновляет данные аккаунта по ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id path int true "ID аккаунта"
// @Param        request body UpdateAccountRequest true "Данные для обновления"
// @Success      200 {object} Account
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /accounts/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.useCase.Update(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to update account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// Delete godoc
// @Summary      Удалить аккаунт
// @Description  Удаляет аккаунт по ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id path int true "ID аккаунта"
// @Success      204 "No Content"
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /accounts/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to delete account", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
