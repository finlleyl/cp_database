package statistics

import (
	"time"

	"github.com/google/uuid"
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

// AccountStatistics represents statistics for an account
type AccountStatistics struct {
	ID             int64     `json:"id" db:"id"`
	AccountID      int64     `json:"account_id" db:"account_id"`
	Period         Period    `json:"period" db:"period"`
	PeriodStart    time.Time `json:"period_start" db:"period_start"`
	PeriodEnd      time.Time `json:"period_end" db:"period_end"`
	TotalTrades    int       `json:"total_trades" db:"total_trades"`
	WinningTrades  int       `json:"winning_trades" db:"winning_trades"`
	LosingTrades   int       `json:"losing_trades" db:"losing_trades"`
	TotalProfit    float64   `json:"total_profit" db:"total_profit"`
	TotalLoss      float64   `json:"total_loss" db:"total_loss"`
	NetProfit      float64   `json:"net_profit" db:"net_profit"`
	ROI            float64   `json:"roi" db:"roi"`
	MaxDrawdownPct float64   `json:"max_drawdown_pct" db:"max_drawdown_pct"`
	WinRate        float64   `json:"win_rate" db:"win_rate"`
	AverageWin     float64   `json:"average_win" db:"average_win"`
	AverageLoss    float64   `json:"average_loss" db:"average_loss"`
	ProfitFactor   float64   `json:"profit_factor" db:"profit_factor"`
	SharpeRatio    float64   `json:"sharpe_ratio" db:"sharpe_ratio"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

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
	UserID           int64            `json:"user_id" db:"user_id"`
	TotalIncome      float64          `json:"total_income" db:"total_income"`
	PerformanceFees  float64          `json:"performance_fees" db:"performance_fees"`
	ManagementFees   float64          `json:"management_fees" db:"management_fees"`
	RegistrationFees float64          `json:"registration_fees" db:"registration_fees"`
	ByStrategy       []StrategyIncome `json:"by_strategy"`
}

// StrategyIncome represents income from a single strategy
type StrategyIncome struct {
	StrategyUUID     uuid.UUID `json:"strategy_uuid" db:"strategy_uuid"`
	StrategyNickname string    `json:"strategy_nickname" db:"strategy_nickname"`
	TotalIncome      float64   `json:"total_income" db:"total_income"`
	PerformanceFees  float64   `json:"performance_fees" db:"performance_fees"`
	ManagementFees   float64   `json:"management_fees" db:"management_fees"`
	Subscribers      int       `json:"subscribers" db:"subscribers"`
}

// Commission represents a commission record
type Commission struct {
	ID               int64          `json:"id" db:"id"`
	SubscriptionUUID uuid.UUID      `json:"subscription_uuid" db:"subscription_uuid"`
	StrategyUUID     uuid.UUID      `json:"strategy_uuid" db:"strategy_uuid"`
	Type             CommissionType `json:"type" db:"type"`
	Amount           float64        `json:"amount" db:"amount"`
	CalculatedFrom   time.Time      `json:"calculated_from" db:"calculated_from"`
	CalculatedTo     time.Time      `json:"calculated_to" db:"calculated_to"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
}

// CommissionType represents the type of commission
type CommissionType string

const (
	CommissionTypePerformance  CommissionType = "performance"
	CommissionTypeManagement   CommissionType = "management"
	CommissionTypeRegistration CommissionType = "registration"
)
