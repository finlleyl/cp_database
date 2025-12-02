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
// @Summary List audit logs
// @Description Get paginated list of audit logs with optional filters
// @Tags Audit
// @Accept json
// @Produce json
// @Param entity_name query string false "Filter by entity name" Enums(users, accounts, strategies, offers, subscriptions, trades)
// @Param entity_pk query string false "Filter by entity primary key"
// @Param operation query string false "Filter by operation" Enums(insert, update, delete)
// @Param changed_by query int false "Filter by user who made the change"
// @Param from query string false "Filter from date (RFC3339)"
// @Param to query string false "Filter to date (RFC3339)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} common.PaginatedResult[AuditLog]
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/audit [get]
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

// GetByEntity handles GET /api/v1/audit/:entity_name/:entity_pk
// @Summary Get audit history for entity
// @Description Get all audit logs for a specific entity (e.g., /api/v1/audit/strategies/123)
// @Tags Audit
// @Accept json
// @Produce json
// @Param entity_name path string true "Entity name" Enums(users, accounts, strategies, offers, subscriptions, trades)
// @Param entity_pk path string true "Entity primary key"
// @Success 200 {array} AuditLog
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/audit/{entity_name}/{entity_pk} [get]
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

// GetStats handles GET /api/v1/audit/stats
// @Summary Get audit statistics
// @Description Get audit statistics grouped by entity and operation
// @Tags Audit
// @Accept json
// @Produce json
// @Param entity_name query string false "Filter by entity name"
// @Param from query string false "Filter from date (RFC3339)"
// @Param to query string false "Filter to date (RFC3339)"
// @Success 200 {array} AuditStats
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/audit/stats [get]
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
