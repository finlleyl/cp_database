package statistics

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	statistics := router.Group("/statistics")
	{
		statistics.GET("/leaderboard", handler.GetStrategyLeaderboard)
		statistics.GET("/investor-portfolio", handler.GetInvestorPortfolio)
		statistics.GET("/master-income", handler.GetMasterIncome)
	}
}
