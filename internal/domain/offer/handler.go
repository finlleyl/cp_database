package offer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for offer operations
type Handler struct {
	useCase UseCase
	logger  *zap.Logger
}

// NewHandler creates a new offer handler
func NewHandler(useCase UseCase, logger *zap.Logger) *Handler {
	return &Handler{useCase: useCase, logger: logger}
}

// Create handles POST /api/v1/offers
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

// GetByUUID handles GET /api/v1/offers/:uuid
func (h *Handler) GetByUUID(c *gin.Context) {
	offerUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offer uuid"})
		return
	}

	offer, err := h.useCase.GetByUUID(c.Request.Context(), offerUUID)
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

// List handles GET /api/v1/offers
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

// Update handles PUT /api/v1/offers/:uuid
func (h *Handler) Update(c *gin.Context) {
	offerUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offer uuid"})
		return
	}

	var req UpdateOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offer, err := h.useCase.Update(c.Request.Context(), offerUUID, &req)
	if err != nil {
		h.logger.Error("Failed to update offer", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, offer)
}

// ChangeStatus handles POST /api/v1/offers/:uuid/status
func (h *Handler) ChangeStatus(c *gin.Context) {
	offerUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offer uuid"})
		return
	}

	var req ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offer, err := h.useCase.ChangeStatus(c.Request.Context(), offerUUID, &req)
	if err != nil {
		h.logger.Error("Failed to change offer status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, offer)
}
