package statistics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for statistics operations
type Handler struct {
	useCase UseCase
	logger  *zap.Logger
}

// NewHandler creates a new statistics handler
func NewHandler(useCase UseCase, logger *zap.Logger) *Handler {
	return &Handler{useCase: useCase, logger: logger}
}

// GetStrategyLeaderboard handles GET /api/v1/statistics/strategies/leaderboard
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

// GetInvestorPortfolio handles GET /api/v1/statistics/investor-portfolio
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

// GetMasterIncome handles GET /api/v1/statistics/master-income
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
