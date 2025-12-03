package strategy

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
// @Summary      Создать стратегию
// @Description  Создаёт новую торговую стратегию
// @Tags         strategies
// @Accept       json
// @Produce      json
// @Param        request body CreateStrategyRequest true "Данные стратегии"
// @Success      201 {object} Strategy
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /strategies [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	strategy, err := h.useCase.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create strategy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, strategy)
}

// GetByID godoc
// @Summary      Получить стратегию по ID
// @Description  Возвращает стратегию по её идентификатору
// @Tags         strategies
// @Accept       json
// @Produce      json
// @Param        id path int true "ID стратегии"
// @Success      200 {object} GetStrategyByIDResponse
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /strategies/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	strategyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid strategy id"})
		return
	}

	strategy, err := h.useCase.GetByID(c.Request.Context(), strategyID)
	if err != nil {
		h.logger.Error("Failed to get strategy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if strategy == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "strategy not found"})
		return
	}

	c.JSON(http.StatusOK, strategy)
}

// List godoc
// @Summary      Список стратегий
// @Description  Возвращает список стратегий с пагинацией и фильтрами
// @Tags         strategies
// @Accept       json
// @Produce      json
// @Param        status query string false "Фильтр по статусу"
// @Param        min_roi query number false "Минимальный ROI"
// @Param        max_drawdown_pct query number false "Максимальная просадка"
// @Param        risk_score query int false "Оценка риска"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество записей на странице" default(20)
// @Success      200 {object} StrategyListResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /strategies [get]
func (h *Handler) List(c *gin.Context) {
	var filter StrategyFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.useCase.List(c.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to list strategies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Update godoc
// @Summary      Обновить стратегию
// @Description  Обновляет данные стратегии по ID
// @Tags         strategies
// @Accept       json
// @Produce      json
// @Param        id path int true "ID стратегии"
// @Param        request body UpdateStrategyRequest true "Данные для обновления"
// @Success      200 {object} Strategy
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /strategies/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	strategyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid strategy id"})
		return
	}

	var req UpdateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	strategy, err := h.useCase.Update(c.Request.Context(), strategyID, &req)
	if err != nil {
		h.logger.Error("Failed to update strategy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, strategy)
}

// ChangeStatus godoc
// @Summary      Изменить статус стратегии
// @Description  Изменяет статус стратегии (active, archived, deleted)
// @Tags         strategies
// @Accept       json
// @Produce      json
// @Param        id path int true "ID стратегии"
// @Param        request body ChangeStatusRequest true "Новый статус"
// @Success      200 {object} Strategy
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /strategies/{id}/status [patch]
func (h *Handler) ChangeStatus(c *gin.Context) {
	strategyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid strategy id"})
		return
	}

	var req ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	strategy, err := h.useCase.ChangeStatus(c.Request.Context(), strategyID, &req)
	if err != nil {
		h.logger.Error("Failed to change strategy status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, strategy)
}

// GetSummary godoc
// @Summary      Получить сводку по стратегии
// @Description  Возвращает суммарную статистику по стратегии
// @Tags         strategies
// @Accept       json
// @Produce      json
// @Param        id path int true "ID стратегии"
// @Success      200 {object} StrategySummary
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /strategies/{id}/summary [get]
func (h *Handler) GetSummary(c *gin.Context) {
	strategyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid strategy id"})
		return
	}

	summary, err := h.useCase.GetSummary(c.Request.Context(), strategyID)
	if err != nil {
		h.logger.Error("Failed to get strategy summary", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
