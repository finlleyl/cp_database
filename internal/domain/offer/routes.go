package offer

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	offers := rg.Group("/offers")
	{
		offers.POST("", h.Create)
		offers.GET("", h.List)
		offers.GET("/:id", h.GetByID)
		offers.PUT("/:id", h.Update)
		offers.POST("/:id/status", h.ChangeStatus)
	}
}
