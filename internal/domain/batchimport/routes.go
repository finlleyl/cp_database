package batchimport

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers batch import routes to the given router group
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	importGroup := rg.Group("/import")
	{
		// POST /api/v1/import/trades - Import trades from file
		importGroup.POST("/trades", h.ImportTrades)

		// GET /api/v1/import - List all import jobs
		importGroup.GET("", h.ListJobs)

		// GET /api/v1/import/:id - Get import job details
		importGroup.GET("/:id", h.GetJobByID)

		// GET /api/v1/import/:id/errors - Get import job errors
		importGroup.GET("/:id/errors", h.GetJobErrors)

		// GET /api/v1/import/:id/summary - Get import job summary
		importGroup.GET("/:id/summary", h.GetJobSummary)
	}
}
