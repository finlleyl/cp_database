package trade

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
// @Summary      Создать сделку
// @Description  Создаёт новую торговую сделку
// @Tags         trades
// @Accept       json
// @Produce      json
// @Param        request body CreateTradeRequest true "Данные сделки"
// @Success      201 {object} Trade
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /trades [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trade, err := h.useCase.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create trade", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, trade)
}

// List godoc
// @Summary      Список сделок
// @Description  Возвращает список сделок с пагинацией и фильтрами
// @Tags         trades
// @Accept       json
// @Produce      json
// @Param        strategy_id query int false "Фильтр по ID стратегии"
// @Param        from query string false "Начало периода (RFC3339)"
// @Param        to query string false "Конец периода (RFC3339)"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество записей на странице" default(20)
// @Success      200 {object} TradeListResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /trades [get]
func (h *Handler) List(c *gin.Context) {
	var filter TradeFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.useCase.List(c.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to list trades", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CopyTrade godoc
// @Summary      Копировать сделку
// @Description  Копирует сделку на указанные подписки
// @Tags         trades
// @Accept       json
// @Produce      json
// @Param        id path int true "ID сделки"
// @Param        request body CopyTradeRequest false "ID подписок для копирования"
// @Success      201 {object} CopyTradeResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /trades/{id}/copy [post]
func (h *Handler) CopyTrade(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	var req CopyTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		req = CopyTradeRequest{}
	}

	copiedTrades, err := h.useCase.CopyTrade(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to copy trade", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"copied_count":  len(copiedTrades),
		"copied_trades": copiedTrades,
	})
}

// ListCopiedTrades godoc
// @Summary      Список скопированных сделок
// @Description  Возвращает список скопированных сделок с фильтрами
// @Tags         trades
// @Accept       json
// @Produce      json
// @Param        subscription_id query int false "Фильтр по ID подписки"
// @Param        trade_id query int false "Фильтр по ID оригинальной сделки"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество записей на странице" default(20)
// @Success      200 {object} CopiedTradeListResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /trades/copied [get]
func (h *Handler) ListCopiedTrades(c *gin.Context) {
	var filter CopiedTradeFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.useCase.ListCopiedTrades(c.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to list copied trades", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
