package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for subscription operations
type Handler struct {
	useCase UseCase
	logger  *zap.Logger
}

// NewHandler creates a new subscription handler
func NewHandler(useCase UseCase, logger *zap.Logger) *Handler {
	return &Handler{useCase: useCase, logger: logger}
}

// Create handles POST /api/v1/subscriptions
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

// GetByUUID handles GET /api/v1/subscriptions/:uuid
func (h *Handler) GetByUUID(c *gin.Context) {
	subscriptionUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription uuid"})
		return
	}

	subscription, err := h.useCase.GetByUUID(c.Request.Context(), subscriptionUUID)
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

// List handles GET /api/v1/subscriptions
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

// Update handles PUT /api/v1/subscriptions/:uuid
func (h *Handler) Update(c *gin.Context) {
	subscriptionUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription uuid"})
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.useCase.Update(c.Request.Context(), subscriptionUUID, &req)
	if err != nil {
		h.logger.Error("Failed to update subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// ChangeStatus handles POST /api/v1/subscriptions/:uuid/status
func (h *Handler) ChangeStatus(c *gin.Context) {
	subscriptionUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription uuid"})
		return
	}

	var req ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get changedBy from auth context
	changedBy := int64(0)

	subscription, err := h.useCase.ChangeStatus(c.Request.Context(), subscriptionUUID, &req, changedBy)
	if err != nil {
		h.logger.Error("Failed to change subscription status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// GetStatusHistory handles GET /api/v1/subscriptions/:uuid/status-history
func (h *Handler) GetStatusHistory(c *gin.Context) {
	subscriptionUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription uuid"})
		return
	}

	history, err := h.useCase.GetStatusHistory(c.Request.Context(), subscriptionUUID)
	if err != nil {
		h.logger.Error("Failed to get subscription status history", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
