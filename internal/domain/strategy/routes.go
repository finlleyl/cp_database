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
		strategies.GET("/:uuid", h.GetByUUID)
		strategies.PUT("/:uuid", h.Update)
		strategies.POST("/:uuid/status", h.ChangeStatus)
	}
}
