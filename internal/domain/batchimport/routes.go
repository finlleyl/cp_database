package batchimport

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers batch import routes to the given router group
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	batchImport := rg.Group("/batch-import")
	{
		batchImport.POST("/trades", h.ImportTrades)
		batchImport.GET("/jobs", h.ListJobs)
		batchImport.GET("/jobs/:id", h.GetJobByID)
		batchImport.GET("/jobs/:id/errors", h.GetJobErrors)
	}
}
