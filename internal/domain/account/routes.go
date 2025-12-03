package account

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	accounts := rg.Group("/accounts")
	{
		accounts.POST("", h.Create)
		accounts.GET("", h.List)
		accounts.GET("/:id", h.GetByID)
		accounts.PUT("/:id", h.Update)
		accounts.DELETE("/:id", h.Delete)
	}
}
