package strategy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for strategy operations
type Handler struct {
	useCase UseCase
	logger  *zap.Logger
}

// NewHandler creates a new strategy handler
func NewHandler(useCase UseCase, logger *zap.Logger) *Handler {
	return &Handler{useCase: useCase, logger: logger}
}

// Create handles POST /api/v1/strategies
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

// GetByUUID handles GET /api/v1/strategies/:uuid
func (h *Handler) GetByUUID(c *gin.Context) {
	strategyUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid strategy uuid"})
		return
	}

	strategy, err := h.useCase.GetByUUID(c.Request.Context(), strategyUUID)
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

// List handles GET /api/v1/strategies
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

// Update handles PUT /api/v1/strategies/:uuid
func (h *Handler) Update(c *gin.Context) {
	strategyUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid strategy uuid"})
		return
	}

	var req UpdateStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	strategy, err := h.useCase.Update(c.Request.Context(), strategyUUID, &req)
	if err != nil {
		h.logger.Error("Failed to update strategy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, strategy)
}

// ChangeStatus handles POST /api/v1/strategies/:uuid/status
func (h *Handler) ChangeStatus(c *gin.Context) {
	strategyUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid strategy uuid"})
		return
	}

	var req ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	strategy, err := h.useCase.ChangeStatus(c.Request.Context(), strategyUUID, &req)
	if err != nil {
		h.logger.Error("Failed to change strategy status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, strategy)
}
