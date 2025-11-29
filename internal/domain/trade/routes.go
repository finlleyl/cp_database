package trade

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers trade routes to the given router group
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	trades := rg.Group("/trades")
	{
		trades.POST("", h.Create)
		trades.GET("", h.List)
		trades.POST("/:id/copy", h.CopyTrade)
	}
	
	// Copied trades endpoint
	rg.GET("/copied-trades", h.ListCopiedTrades)
}

