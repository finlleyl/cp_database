package batchimport

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for batch import operations
type Handler struct {
	useCase UseCase
	logger  *zap.Logger
}

// NewHandler creates a new batch import handler
func NewHandler(useCase UseCase, logger *zap.Logger) *Handler {
	return &Handler{useCase: useCase, logger: logger}
}

// ImportTrades handles POST /api/v1/batch-import/trades
func (h *Handler) ImportTrades(c *gin.Context) {
	var req ImportTradesRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	job, err := h.useCase.ImportTrades(c.Request.Context(), &req, file, header.Filename, header.Size)
	if err != nil {
		h.logger.Error("Failed to import trades", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, job)
}

// ListJobs handles GET /api/v1/batch-import/jobs
func (h *Handler) ListJobs(c *gin.Context) {
	var filter JobFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.useCase.ListJobs(c.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to list import jobs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetJobByID handles GET /api/v1/batch-import/jobs/:id
func (h *Handler) GetJobByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	job, err := h.useCase.GetJobByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get import job", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if job == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// GetJobErrors handles GET /api/v1/batch-import/jobs/:id/errors
func (h *Handler) GetJobErrors(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	errors, err := h.useCase.GetJobErrors(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get import job errors", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, errors)
}

