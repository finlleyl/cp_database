package audit

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	auditGroup := rg.Group("/audit")
	{

		auditGroup.GET("", h.List)

		auditGroup.GET("/stats", h.GetStats)

		auditGroup.GET("/:entity_name/:entity_pk", h.GetByEntity)
	}
}
