package audit

import (
	"net/http"

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

// List godoc
// @Summary      Список аудит-логов
// @Description  Возвращает список записей аудита с фильтрами
// @Tags         audit
// @Accept       json
// @Produce      json
// @Param        entity_name query string false "Фильтр по имени сущности (users/accounts/strategies/offers/subscriptions/trades)"
// @Param        entity_pk query string false "Фильтр по первичному ключу сущности"
// @Param        operation query string false "Фильтр по операции (insert/update/delete)"
// @Param        changed_by query int false "Фильтр по ID пользователя"
// @Param        from query string false "Начало периода (RFC3339)"
// @Param        to query string false "Конец периода (RFC3339)"
// @Param        page query int false "Номер страницы" default(1)
// @Param        limit query int false "Количество записей на странице" default(20)
// @Success      200 {object} AuditListResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /audit [get]
func (h *Handler) List(c *gin.Context) {
	var filter AuditFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.useCase.List(c.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to list audit logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetByEntity godoc
// @Summary      История изменений сущности
// @Description  Возвращает историю изменений для конкретной сущности
// @Tags         audit
// @Accept       json
// @Produce      json
// @Param        entity_name path string true "Имя сущности (users/accounts/strategies/offers/subscriptions/trades)"
// @Param        entity_pk path string true "Первичный ключ сущности"
// @Success      200 {array} AuditLog
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /audit/{entity_name}/{entity_pk} [get]
func (h *Handler) GetByEntity(c *gin.Context) {
	entityName := c.Param("entity_name")
	entityPK := c.Param("entity_pk")

	if entityName == "" || entityPK == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_name and entity_pk are required"})
		return
	}

	logs, err := h.useCase.GetByEntity(c.Request.Context(), entityName, entityPK)
	if err != nil {
		h.logger.Error("Failed to get audit logs by entity",
			zap.String("entity_name", entityName),
			zap.String("entity_pk", entityPK),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetStats godoc
// @Summary      Статистика аудита
// @Description  Возвращает агрегированную статистику по аудит-логам
// @Tags         audit
// @Accept       json
// @Produce      json
// @Param        entity_name query string false "Фильтр по имени сущности"
// @Param        from query string false "Начало периода (RFC3339)"
// @Param        to query string false "Конец периода (RFC3339)"
// @Success      200 {array} AuditStats
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /audit/stats [get]
func (h *Handler) GetStats(c *gin.Context) {
	var filter AuditStatsFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := h.useCase.GetStats(c.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to get audit stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
