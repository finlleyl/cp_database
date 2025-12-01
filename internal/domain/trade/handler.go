package trade

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for trade operations
type Handler struct {
	useCase UseCase
	logger  *zap.Logger
}

// NewHandler creates a new trade handler
func NewHandler(useCase UseCase, logger *zap.Logger) *Handler {
	return &Handler{useCase: useCase, logger: logger}
}

// Create handles POST /api/v1/trades
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

// List handles GET /api/v1/trades
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

// CopyTrade handles POST /api/v1/trades/:id/copy
func (h *Handler) CopyTrade(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	var req CopyTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Allow empty body for copying to all subscriptions
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

// ListCopiedTrades handles GET /api/v1/copied-trades
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
