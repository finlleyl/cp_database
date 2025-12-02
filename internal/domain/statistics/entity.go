package statistics

import (
	"time"
)

// Period represents a statistics period
type Period string

const (
	PeriodDay   Period = "day"
	PeriodWeek  Period = "week"
	PeriodMonth Period = "month"
	PeriodYear  Period = "year"
	PeriodAll   Period = "all"
)

// StrategyLeaderboard represents a strategy in the leaderboard
// Matches fn_get_strategy_leaderboard return type
type StrategyLeaderboard struct {
	StrategyID          int64   `json:"strategy_id" db:"strategy_id"`
	Title               string  `json:"title" db:"title"`
	TotalProfit         float64 `json:"total_profit" db:"total_profit"`
	TotalCommissions    float64 `json:"total_commissions" db:"total_commissions"`
	ActiveSubscriptions int     `json:"active_subscriptions" db:"active_subscriptions"`
}

// InvestorPortfolio represents an investor's portfolio statistics
// Response wrapper for investor portfolio endpoint
type InvestorPortfolio struct {
	UserID        int64           `json:"user_id"`
	Subscriptions []PortfolioItem `json:"subscriptions"`
}

// PortfolioItem represents a single subscription in the portfolio
// Matches fn_get_investor_portfolio return type
type PortfolioItem struct {
	SubscriptionID    int64   `json:"subscription_id" db:"subscription_id"`
	StrategyID        int64   `json:"strategy_id" db:"strategy_id"`
	StrategyTitle     string  `json:"strategy_title" db:"strategy_title"`
	TotalProfit       float64 `json:"total_profit" db:"total_profit"`
	CopiedTradesCount int64   `json:"copied_trades_count" db:"copied_trades_count"`
}

// MasterIncome represents a master's income from commissions
type MasterIncome struct {
	UserID           int64   `json:"user_id"`
	TotalIncome      float64 `json:"total_income"`
	PerformanceFees  float64 `json:"performance_fees"`
	ManagementFees   float64 `json:"management_fees"`
	RegistrationFees float64 `json:"registration_fees"`
}

// Commission represents a commission record
type Commission struct {
	ID             int64          `json:"id" db:"id"`
	SubscriptionID int64          `json:"subscription_id" db:"subscription_id"`
	Type           CommissionType `json:"type" db:"type"`
	Amount         float64        `json:"amount" db:"amount"`
	PeriodFrom     *time.Time     `json:"period_from,omitempty" db:"period_from"`
	PeriodTo       *time.Time     `json:"period_to,omitempty" db:"period_to"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
}

// CommissionType represents the type of commission
type CommissionType string

const (
	CommissionTypePerformance  CommissionType = "performance"
	CommissionTypeManagement   CommissionType = "management"
	CommissionTypeRegistration CommissionType = "registration"
)

// CreateCommissionRequest represents the request to create a commission
type CreateCommissionRequest struct {
	SubscriptionID int64          `json:"subscription_id" binding:"required"`
	Type           CommissionType `json:"type" binding:"required,oneof=performance management registration"`
	Amount         float64        `json:"amount" binding:"required,gte=0"`
	PeriodFrom     *time.Time     `json:"period_from,omitempty"`
	PeriodTo       *time.Time     `json:"period_to,omitempty"`
}
