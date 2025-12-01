package subscription

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers subscription routes to the given router group
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	subscriptions := rg.Group("/subscriptions")
	{
		subscriptions.POST("", h.Create)
		subscriptions.GET("", h.List)
		subscriptions.GET("/:uuid", h.GetByUUID)
		subscriptions.PUT("/:uuid", h.Update)
		subscriptions.POST("/:uuid/status", h.ChangeStatus)
		subscriptions.GET("/:uuid/status-history", h.GetStatusHistory)
	}
}
