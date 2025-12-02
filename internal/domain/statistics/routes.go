package statistics

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all statistics routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	statistics := router.Group("/statistics")
	{
		statistics.GET("/strategies/leaderboard", handler.GetStrategyLeaderboard)
		statistics.GET("/investor-portfolio", handler.GetInvestorPortfolio)
		statistics.GET("/master-income", handler.GetMasterIncome)
	}
}
