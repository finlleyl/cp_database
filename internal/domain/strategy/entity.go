package strategy

import (
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/google/uuid"
)

// Strategy represents a master's trading strategy
type Strategy struct {
	UUID             uuid.UUID             `json:"uuid" db:"uuid"`
	AccountID        int64                 `json:"account_id" db:"account_id"`
	UserID           int64                 `json:"user_id" db:"user_id"`
	Nickname         string                `json:"nickname" db:"nickname"`
	Summary          string                `json:"summary" db:"summary"`
	PaymentAccountID *int64                `json:"payment_account_id,omitempty" db:"payment_account_id"`
	AvatarURL        string                `json:"avatar_url" db:"avatar_url"`
	Status           common.StrategyStatus `json:"status" db:"status"`
	StatusReason     string                `json:"status_reason" db:"status_reason"`
	ROI              float64               `json:"roi" db:"roi"`
	MaxDrawdownPct   float64               `json:"max_drawdown_pct" db:"max_drawdown_pct"`
	RiskScore        int                   `json:"risk_score" db:"risk_score"`
	CreatedAt        time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at" db:"updated_at"`
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
