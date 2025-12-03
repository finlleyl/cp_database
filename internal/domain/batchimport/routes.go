package batchimport

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	importGroup := rg.Group("/import")
	{

		importGroup.POST("/trades", h.ImportTrades)

		importGroup.GET("", h.ListJobs)

		importGroup.GET("/:id", h.GetJobByID)

		importGroup.GET("/:id/errors", h.GetJobErrors)

		importGroup.GET("/:id/summary", h.GetJobSummary)
	}
}
