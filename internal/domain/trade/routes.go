package trade

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	trades := rg.Group("/trades")
	{
		trades.POST("", h.Create)
		trades.GET("", h.List)
		trades.POST("/:id/copy", h.CopyTrade)
	}

	rg.GET("/copied-trades", h.ListCopiedTrades)
}
