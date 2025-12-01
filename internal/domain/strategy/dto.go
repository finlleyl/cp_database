package strategy

import (
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

// GetStrategyByIDResponse represents the response for getting strategy by ID
// Contains aggregated statistics from vw_strategy_performance view
type GetStrategyByIDResponse struct {
	ID                  int64                 `json:"id" db:"id"`
	Title               string                `json:"title" db:"title"`
	Status              common.StrategyStatus `json:"status" db:"status"`
	TotalSubscriptions  int64                 `json:"total_subscriptions" db:"total_subscriptions"`
	ActiveSubscriptions int64                 `json:"active_subscriptions" db:"active_subscriptions"`
	TotalCopiedTrades   int64                 `json:"total_copied_trades" db:"total_copied_trades"`
	TotalProfit         float64               `json:"total_profit" db:"total_profit"`
	TotalCommissions    float64               `json:"total_commissions" db:"total_commissions"`
	UpdatedAt           time.Time             `json:"updated_at" db:"updated_at"`
}

// CreateStrategyRequest represents the request to create a new strategy
type CreateStrategyRequest struct {
	AccountID        int64  `json:"account_id" binding:"required"`
	UserID           int64  `json:"user_id" binding:"required"`
	Nickname         string `json:"nickname" binding:"required"`
	Summary          string `json:"summary"`
	PaymentAccountID *int64 `json:"payment_account_id,omitempty"`
	AvatarURL        string `json:"avatar_url"`
}

// UpdateStrategyRequest represents the request to update a strategy
type UpdateStrategyRequest struct {
	Nickname         *string `json:"nickname,omitempty"`
	Summary          *string `json:"summary,omitempty"`
	PaymentAccountID *int64  `json:"payment_account_id,omitempty"`
	AvatarURL        *string `json:"avatar_url,omitempty"`
}

// ChangeStatusRequest represents the request to change strategy status
type ChangeStatusRequest struct {
	Status       common.StrategyStatus `json:"status" binding:"required,oneof=active archived deleted"`
	StatusReason string                `json:"status_reason"`
}

// StrategyFilter represents filter parameters for strategy search
type StrategyFilter struct {
	Status         common.StrategyStatus `form:"status"`
	MinROI         *float64              `form:"min_roi"`
	MaxDrawdownPct *float64              `form:"max_drawdown_pct"`
	RiskScore      *int                  `form:"risk_score"`
	Search         string                `form:"search"`
	common.Pagination
}
