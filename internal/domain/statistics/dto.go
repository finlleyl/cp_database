package statistics

import (
	"github.com/finlleyl/cp_database/internal/domain/common"
)

type LeaderboardRequest struct {
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

type InvestorPortfolioRequest struct {
	UserID int64 `form:"user_id" binding:"required"`
}

type MasterIncomeRequest struct {
	UserID int64 `form:"user_id" binding:"required"`
	common.TimeRange
}

type AccountStatisticsRequest struct {
	AccountID int64  `uri:"account_id" binding:"required"`
	Period    Period `form:"period" binding:"omitempty,oneof=day week month year all"`
}
