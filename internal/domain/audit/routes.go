package audit

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers audit routes to the given router group
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/audit", h.List)
}
