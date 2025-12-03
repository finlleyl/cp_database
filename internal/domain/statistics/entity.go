package statistics

import (
	"time"
)

type Period string

const (
	PeriodDay   Period = "day"
	PeriodWeek  Period = "week"
	PeriodMonth Period = "month"
	PeriodYear  Period = "year"
	PeriodAll   Period = "all"
)

type StrategyLeaderboard struct {
	StrategyID          int64   `json:"strategy_id" db:"strategy_id"`
	Title               string  `json:"title" db:"title"`
	TotalProfit         float64 `json:"total_profit" db:"total_profit"`
	TotalCommissions    float64 `json:"total_commissions" db:"total_commissions"`
	ActiveSubscriptions int     `json:"active_subscriptions" db:"active_subscriptions"`
}

type InvestorPortfolio struct {
	UserID        int64           `json:"user_id"`
	Subscriptions []PortfolioItem `json:"subscriptions"`
}

type PortfolioItem struct {
	SubscriptionID    int64   `json:"subscription_id" db:"subscription_id"`
	StrategyID        int64   `json:"strategy_id" db:"strategy_id"`
	StrategyTitle     string  `json:"strategy_title" db:"strategy_title"`
	TotalProfit       float64 `json:"total_profit" db:"total_profit"`
	CopiedTradesCount int64   `json:"copied_trades_count" db:"copied_trades_count"`
}

type MasterIncome struct {
	UserID           int64   `json:"user_id"`
	TotalIncome      float64 `json:"total_income"`
	PerformanceFees  float64 `json:"performance_fees"`
	ManagementFees   float64 `json:"management_fees"`
	RegistrationFees float64 `json:"registration_fees"`
}

type Commission struct {
	ID             int64          `json:"id" db:"id"`
	SubscriptionID int64          `json:"subscription_id" db:"subscription_id"`
	Type           CommissionType `json:"type" db:"type"`
	Amount         float64        `json:"amount" db:"amount"`
	PeriodFrom     *time.Time     `json:"period_from,omitempty" db:"period_from"`
	PeriodTo       *time.Time     `json:"period_to,omitempty" db:"period_to"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
}

type CommissionType string

const (
	CommissionTypePerformance  CommissionType = "performance"
	CommissionTypeManagement   CommissionType = "management"
	CommissionTypeRegistration CommissionType = "registration"
)

type CreateCommissionRequest struct {
	SubscriptionID int64          `json:"subscription_id" binding:"required"`
	Type           CommissionType `json:"type" binding:"required,oneof=performance management registration"`
	Amount         float64        `json:"amount" binding:"required,gte=0"`
	PeriodFrom     *time.Time     `json:"period_from,omitempty"`
	PeriodTo       *time.Time     `json:"period_to,omitempty"`
}
