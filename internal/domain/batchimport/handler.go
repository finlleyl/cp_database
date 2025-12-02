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

// ImportTrades handles POST /api/v1/import/trades
// @Summary Import trades from file
// @Description Upload a CSV or JSON file to import trades for a strategy
// @Tags Import
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Trade data file (CSV or JSON)"
// @Param strategy_uuid formData string true "Strategy UUID"
// @Param account_id formData int true "Account ID"
// @Param file_format formData string true "File format" Enums(csv, json)
// @Success 202 {object} ImportJob
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/import/trades [post]
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

	job, err := h.useCase.ImportTrades(c.Request.Context(), &req, file, header.Filename)
	if err != nil {
		h.logger.Error("Failed to import trades", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, job)
}

// ListJobs handles GET /api/v1/import
// @Summary List import jobs
// @Description Get paginated list of import jobs with optional filters
// @Tags Import
// @Accept json
// @Produce json
// @Param type query string false "Filter by job type" Enums(trades, accounts, statistics)
// @Param status query string false "Filter by status" Enums(pending, running, success, failed)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} common.PaginatedResult[ImportJob]
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/import [get]
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

// GetJobByID handles GET /api/v1/import/:id
// @Summary Get import job by ID
// @Description Get details of a specific import job
// @Tags Import
// @Accept json
// @Produce json
// @Param id path int true "Job ID"
// @Success 200 {object} ImportJob
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/import/{id} [get]
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

// GetJobErrors handles GET /api/v1/import/:id/errors
// @Summary Get import job errors
// @Description Get paginated list of errors for a specific import job
// @Tags Import
// @Accept json
// @Produce json
// @Param id path int true "Job ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} common.PaginatedResult[ImportJobError]
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/import/{id}/errors [get]
func (h *Handler) GetJobErrors(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	var filter ErrorFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errors, err := h.useCase.GetJobErrors(c.Request.Context(), id, &filter)
	if err != nil {
		h.logger.Error("Failed to get import job errors", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, errors)
}

// GetJobSummary handles GET /api/v1/import/:id/summary
// @Summary Get import job summary
// @Description Get summary statistics for a completed import job
// @Tags Import
// @Accept json
// @Produce json
// @Param id path int true "Job ID"
// @Success 200 {object} ImportJobSummary
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/import/{id}/summary [get]
func (h *Handler) GetJobSummary(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	summary, err := h.useCase.GetJobSummary(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get import job summary", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
