package subscription

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
// @Summary      Создать подписку
// @Description  Создаёт новую подписку инвестора на оффер
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request body CreateSubscriptionRequest true "Данные подписки"
// @Success      201 {object} Subscription
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.useCase.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// GetByUUID godoc
// @Summary      Получить подписку по ID
// @Description  Возвращает подписку по её идентификатору
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path int true "ID подписки"
// @Success      200 {object} Subscription
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/{id} [get]
func (h *Handler) GetByUUID(c *gin.Context) {
	subscriptionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription uuid"})
		return
	}

	subscription, err := h.useCase.GetByID(c.Request.Context(), subscriptionID)
	if err != nil {
		h.logger.Error("Failed to get subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if subscription == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// List godoc
// @Summary      Список подписок
// @Description  Возвращает список подписок с пагинацией и фильтрами
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        user_id query int false "Фильтр по ID пользователя"
// @Param        offer_id query int false "Фильтр по ID оффера"
// @Param        status query string false "Фильтр по статусу"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество записей на странице" default(20)
// @Success      200 {object} SubscriptionListResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions [get]
func (h *Handler) List(c *gin.Context) {
	var filter SubscriptionFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.useCase.List(c.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to list subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Update godoc
// @Summary      Обновить подписку
// @Description  Обновляет данные подписки по ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path int true "ID подписки"
// @Param        request body UpdateSubscriptionRequest true "Данные для обновления"
// @Success      200 {object} Subscription
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	subscriptionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription uuid"})
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.useCase.Update(c.Request.Context(), subscriptionID, &req)
	if err != nil {
		h.logger.Error("Failed to update subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// ChangeStatus godoc
// @Summary      Изменить статус подписки
// @Description  Изменяет статус подписки (active, archived, suspended, deleted)
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path int true "ID подписки"
// @Param        request body ChangeStatusRequest true "Новый статус"
// @Success      200 {object} Subscription
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/{id}/status [patch]
func (h *Handler) ChangeStatus(c *gin.Context) {
	subscriptionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription uuid"})
		return
	}

	var req ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	changedBy := int64(0)

	subscription, err := h.useCase.ChangeStatus(c.Request.Context(), subscriptionID, &req, changedBy)
	if err != nil {
		h.logger.Error("Failed to change subscription status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// GetStatusHistory godoc
// @Summary      История статусов подписки
// @Description  Возвращает историю изменений статуса подписки
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path int true "ID подписки"
// @Success      200 {array} SubscriptionStatusHistory
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/{id}/history [get]
func (h *Handler) GetStatusHistory(c *gin.Context) {
	subscriptionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription uuid"})
		return
	}

	history, err := h.useCase.GetStatusHistory(c.Request.Context(), subscriptionID)
	if err != nil {
		h.logger.Error("Failed to get subscription status history", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
