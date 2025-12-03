package offer

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
// @Summary      Создать оффер
// @Description  Создаёт новый оффер для стратегии
// @Tags         offers
// @Accept       json
// @Produce      json
// @Param        request body CreateOfferRequest true "Данные оффера"
// @Success      201 {object} Offer
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /offers [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offer, err := h.useCase.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create offer", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, offer)
}

// GetByID godoc
// @Summary      Получить оффер по ID
// @Description  Возвращает оффер по его идентификатору
// @Tags         offers
// @Accept       json
// @Produce      json
// @Param        id path int true "ID оффера"
// @Success      200 {object} Offer
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /offers/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	offerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offer id"})
		return
	}

	offer, err := h.useCase.GetByID(c.Request.Context(), offerID)
	if err != nil {
		h.logger.Error("Failed to get offer", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if offer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "offer not found"})
		return
	}

	c.JSON(http.StatusOK, offer)
}

// List godoc
// @Summary      Список офферов
// @Description  Возвращает список офферов с пагинацией и фильтрами
// @Tags         offers
// @Accept       json
// @Produce      json
// @Param        strategy_id query int false "Фильтр по ID стратегии"
// @Param        status query string false "Фильтр по статусу (active/archived/deleted)"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество записей на странице" default(20)
// @Success      200 {object} OfferListResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /offers [get]
func (h *Handler) List(c *gin.Context) {
	var filter OfferFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.useCase.List(c.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to list offers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Update godoc
// @Summary      Обновить оффер
// @Description  Обновляет данные оффера по ID
// @Tags         offers
// @Accept       json
// @Produce      json
// @Param        id path int true "ID оффера"
// @Param        request body UpdateOfferRequest true "Данные для обновления"
// @Success      200 {object} Offer
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /offers/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	offerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offer id"})
		return
	}

	var req UpdateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offer, err := h.useCase.Update(c.Request.Context(), offerID, &req)
	if err != nil {
		h.logger.Error("Failed to update offer", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, offer)
}

// ChangeStatus godoc
// @Summary      Изменить статус оффера
// @Description  Изменяет статус оффера (active, archived, deleted)
// @Tags         offers
// @Accept       json
// @Produce      json
// @Param        id path int true "ID оффера"
// @Param        request body ChangeStatusRequest true "Новый статус"
// @Success      200 {object} Offer
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /offers/{id}/status [patch]
func (h *Handler) ChangeStatus(c *gin.Context) {
	offerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offer id"})
		return
	}

	var req ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offer, err := h.useCase.ChangeStatus(c.Request.Context(), offerID, &req)
	if err != nil {
		h.logger.Error("Failed to change offer status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, offer)
}
