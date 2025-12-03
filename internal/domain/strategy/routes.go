package strategy

import (
	"github.com/gin-gonic/gin"
)

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
