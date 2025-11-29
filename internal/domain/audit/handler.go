package audit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for audit log operations
type Handler struct {
	useCase UseCase
	logger  *zap.Logger
}

// NewHandler creates a new audit handler
func NewHandler(useCase UseCase, logger *zap.Logger) *Handler {
	return &Handler{useCase: useCase, logger: logger}
}

// List handles GET /api/v1/audit
func (h *Handler) List(c *gin.Context) {
	var filter AuditFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.useCase.List(c.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to list audit logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

