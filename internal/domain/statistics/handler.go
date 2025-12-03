package statistics

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

// GetStrategyLeaderboard godoc
// @Summary      Лидерборд стратегий
// @Description  Возвращает топ стратегий по доходности
// @Tags         statistics
// @Accept       json
// @Produce      json
// @Param        limit query int false "Количество записей" default(10)
// @Success      200 {array} LeaderboardEntry
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /statistics/leaderboard [get]
func (h *Handler) GetStrategyLeaderboard(c *gin.Context) {
	var req LeaderboardRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leaderboard, err := h.useCase.GetStrategyLeaderboard(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to get strategy leaderboard", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

// GetInvestorPortfolio godoc
// @Summary      Портфель инвестора
// @Description  Возвращает портфель инвестора с его подписками и статистикой
// @Tags         statistics
// @Accept       json
// @Produce      json
// @Param        user_id query int true "ID пользователя-инвестора"
// @Success      200 {array} PortfolioEntry
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /statistics/portfolio [get]
func (h *Handler) GetInvestorPortfolio(c *gin.Context) {
	var req InvestorPortfolioRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	portfolio, err := h.useCase.GetInvestorPortfolio(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to get investor portfolio", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// GetMasterIncome godoc
// @Summary      Доход мастера
// @Description  Возвращает информацию о доходах мастер-трейдера
// @Tags         statistics
// @Accept       json
// @Produce      json
// @Param        user_id query int true "ID пользователя-мастера"
// @Param        from query string false "Начало периода (RFC3339)"
// @Param        to query string false "Конец периода (RFC3339)"
// @Success      200 {object} MasterIncome
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /statistics/master-income [get]
func (h *Handler) GetMasterIncome(c *gin.Context) {
	var req MasterIncomeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	income, err := h.useCase.GetMasterIncome(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to get master income", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, income)
}
