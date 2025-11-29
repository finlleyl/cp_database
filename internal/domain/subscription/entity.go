package subscription

import (
	"encoding/json"
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/google/uuid"
)

// Subscription represents an investor's subscription to a strategy
type Subscription struct {
	UUID              uuid.UUID                 `json:"uuid" db:"uuid"`
	InvestorAccountID int64                     `json:"investor_account_id" db:"investor_account_id"`
	OfferUUID         uuid.UUID                 `json:"offer_uuid" db:"offer_uuid"`
	UserID            int64                     `json:"user_id" db:"user_id"`
	Config            json.RawMessage           `json:"config" db:"config"`
	RiskRules         json.RawMessage           `json:"risk_rules" db:"risk_rules"`
	Filter            json.RawMessage           `json:"filter" db:"filter"`
	Status            common.SubscriptionStatus `json:"status" db:"status"`
	StatusReason      string                    `json:"status_reason" db:"status_reason"`
	CreatedAt         time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at" db:"updated_at"`
}

// SubscriptionConfig represents the configuration for a subscription
type SubscriptionConfig struct {
	CopyRatio        float64 `json:"copy_ratio"`
	MaxPositionSize  float64 `json:"max_position_size"`
	InvertSignals    bool    `json:"invert_signals"`
	CopyStopLoss     bool    `json:"copy_stop_loss"`
	CopyTakeProfit   bool    `json:"copy_take_profit"`
}

// RiskRules represents the risk rules for a subscription
type RiskRules struct {
	MaxDrawdownPct   float64 `json:"max_drawdown_pct"`
	MaxDailyLoss     float64 `json:"max_daily_loss"`
	StopOnDrawdown   bool    `json:"stop_on_drawdown"`
}

// SubscriptionFilter represents the filter criteria for a subscription
type SubscriptionFilterConfig struct {
	AllowedSymbols   []string `json:"allowed_symbols"`
	BlockedSymbols   []string `json:"blocked_symbols"`
	MinLotSize       float64  `json:"min_lot_size"`
	MaxLotSize       float64  `json:"max_lot_size"`
}

// CreateSubscriptionRequest represents the request to create a new subscription
type CreateSubscriptionRequest struct {
	InvestorAccountID int64           `json:"investor_account_id" binding:"required"`
	OfferUUID         uuid.UUID       `json:"offer_uuid" binding:"required"`
	UserID            int64           `json:"user_id" binding:"required"`
	Config            json.RawMessage `json:"config"`
	RiskRules         json.RawMessage `json:"risk_rules"`
	Filter            json.RawMessage `json:"filter"`
}

// UpdateSubscriptionRequest represents the request to update a subscription
type UpdateSubscriptionRequest struct {
	Config    json.RawMessage `json:"config,omitempty"`
	RiskRules json.RawMessage `json:"risk_rules,omitempty"`
	Filter    json.RawMessage `json:"filter,omitempty"`
}

// ChangeStatusRequest represents the request to change subscription status
type ChangeStatusRequest struct {
	Status       common.SubscriptionStatus `json:"status" binding:"required,oneof=active archived suspended deleted"`
	StatusReason string                    `json:"status_reason"`
}

// SubscriptionStatusHistory represents a status change history record
type SubscriptionStatusHistory struct {
	ID               int64                     `json:"id" db:"id"`
	SubscriptionUUID uuid.UUID                 `json:"subscription_uuid" db:"subscription_uuid"`
	OldStatus        common.SubscriptionStatus `json:"old_status" db:"old_status"`
	NewStatus        common.SubscriptionStatus `json:"new_status" db:"new_status"`
	Reason           string                    `json:"reason" db:"reason"`
	ChangedBy        int64                     `json:"changed_by" db:"changed_by"`
	CreatedAt        time.Time                 `json:"created_at" db:"created_at"`
}

// SubscriptionFilter represents filter parameters for subscription search
type SubscriptionFilter struct {
	UserID   int64                     `form:"user_id"`
	Status   common.SubscriptionStatus `form:"status"`
	common.Pagination
}

