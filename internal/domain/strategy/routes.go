package strategy

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers strategy routes to the given router group
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	strategies := rg.Group("/strategies")
	{
		strategies.POST("", h.Create)
		strategies.GET("", h.List)
		strategies.GET("/:id", h.GetByID)
		strategies.GET("/:id/summary", h.GetSummary)
		strategies.PUT("/:id", h.Update)
		strategies.POST("/:id/status", h.ChangeStatus)
	}
}
