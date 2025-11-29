package offer

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers offer routes to the given router group
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	offers := rg.Group("/offers")
	{
		offers.POST("", h.Create)
		offers.GET("", h.List)
		offers.GET("/:uuid", h.GetByUUID)
		offers.PUT("/:uuid", h.Update)
		offers.POST("/:uuid/status", h.ChangeStatus)
	}
}

