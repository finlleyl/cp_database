package audit

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers audit routes to the given router group
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	auditGroup := rg.Group("/audit")
	{
		// GET /api/v1/audit - List all audit logs with filters
		auditGroup.GET("", h.List)

		// GET /api/v1/audit/stats - Get audit statistics
		auditGroup.GET("/stats", h.GetStats)

		// GET /api/v1/audit/:entity_name/:entity_pk - Get audit history for specific entity
		auditGroup.GET("/:entity_name/:entity_pk", h.GetByEntity)
	}
}
