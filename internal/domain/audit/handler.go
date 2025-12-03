package audit

import (
	"net/http"

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

func (h *Handler) GetByEntity(c *gin.Context) {
	entityName := c.Param("entity_name")
	entityPK := c.Param("entity_pk")

	if entityName == "" || entityPK == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_name and entity_pk are required"})
		return
	}

	logs, err := h.useCase.GetByEntity(c.Request.Context(), entityName, entityPK)
	if err != nil {
		h.logger.Error("Failed to get audit logs by entity",
			zap.String("entity_name", entityName),
			zap.String("entity_pk", entityPK),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func (h *Handler) GetStats(c *gin.Context) {
	var filter AuditStatsFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := h.useCase.GetStats(c.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to get audit stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
