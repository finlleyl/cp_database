package statistics

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers statistics routes to the given router group
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	stats := rg.Group("/statistics")
	{
		stats.GET("/strategies/leaderboard", h.GetStrategyLeaderboard)
		stats.GET("/investor-portfolio", h.GetInvestorPortfolio)
		stats.GET("/master-income", h.GetMasterIncome)
		stats.GET("/accounts/:account_id", h.GetAccountStatistics)
	}
}
