package batchimport

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	useCase UseCase
	logger  *zap.Logger
}

func NewHandler(useCase UseCase, logger *zap.Logger) *Handler {
	return &Handler{useCase: useCase, logger: logger}
}

// ImportTrades godoc
// @Summary      Импортировать сделки
// @Description  Загружает файл со сделками для импорта
// @Tags         import
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Файл со сделками (CSV или JSON)"
// @Param        strategy_id formData int true "ID стратегии"
// @Param        account_id formData int true "ID аккаунта"
// @Param        file_format formData string true "Формат файла (csv/json)"
// @Success      202 {object} ImportJob
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /import/trades [post]
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

// ListJobs godoc
// @Summary      Список задач импорта
// @Description  Возвращает список задач импорта с фильтрами
// @Tags         import
// @Accept       json
// @Produce      json
// @Param        type query string false "Фильтр по типу (trades/accounts/statistics)"
// @Param        status query string false "Фильтр по статусу"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество записей на странице" default(20)
// @Success      200 {object} JobListResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /import/jobs [get]
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

// GetJobByID godoc
// @Summary      Получить задачу импорта по ID
// @Description  Возвращает информацию о задаче импорта
// @Tags         import
// @Accept       json
// @Produce      json
// @Param        id path int true "ID задачи"
// @Success      200 {object} ImportJob
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /import/jobs/{id} [get]
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

// GetJobErrors godoc
// @Summary      Ошибки задачи импорта
// @Description  Возвращает список ошибок для задачи импорта
// @Tags         import
// @Accept       json
// @Produce      json
// @Param        id path int true "ID задачи"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество записей на странице" default(20)
// @Success      200 {object} JobErrorListResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /import/jobs/{id}/errors [get]
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

// GetJobSummary godoc
// @Summary      Сводка по задаче импорта
// @Description  Возвращает сводную информацию о выполнении задачи импорта
// @Tags         import
// @Accept       json
// @Produce      json
// @Param        id path int true "ID задачи"
// @Success      200 {object} ImportJobSummary
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /import/jobs/{id}/summary [get]
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
