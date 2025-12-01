package statistics

import (
	"github.com/finlleyl/cp_database/internal/domain/common"
)

// LeaderboardRequest represents the request for strategy leaderboard
type LeaderboardRequest struct {
	Period Period `form:"period" binding:"omitempty,oneof=day week month year all"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// InvestorPortfolioRequest represents the request for investor portfolio
type InvestorPortfolioRequest struct {
	UserID int64 `form:"user_id" binding:"required"`
}

// MasterIncomeRequest represents the request for master income
type MasterIncomeRequest struct {
	UserID int64 `form:"user_id" binding:"required"`
	common.TimeRange
}

// AccountStatisticsRequest represents the request for account statistics
type AccountStatisticsRequest struct {
	AccountID int64  `uri:"account_id" binding:"required"`
	Period    Period `form:"period" binding:"omitempty,oneof=day week month year all"`
}
