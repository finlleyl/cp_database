package trade

import (
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/google/uuid"
)

// TradeType represents the type of trade
type TradeType string

const (
	TradeTypeBuy  TradeType = "buy"
	TradeTypeSell TradeType = "sell"
)

// TradeStatus represents the status of a trade
type TradeStatus string

const (
	TradeStatusOpen   TradeStatus = "open"
	TradeStatusClosed TradeStatus = "closed"
)

// Trade represents a master's trade
type Trade struct {
	ID           int64       `json:"id" db:"id"`
	StrategyUUID uuid.UUID   `json:"strategy_uuid" db:"strategy_uuid"`
	AccountID    int64       `json:"account_id" db:"account_id"`
	Symbol       string      `json:"symbol" db:"symbol"`
	Type         TradeType   `json:"type" db:"type"`
	Volume       float64     `json:"volume" db:"volume"`
	OpenPrice    float64     `json:"open_price" db:"open_price"`
	ClosePrice   *float64    `json:"close_price,omitempty" db:"close_price"`
	StopLoss     *float64    `json:"stop_loss,omitempty" db:"stop_loss"`
	TakeProfit   *float64    `json:"take_profit,omitempty" db:"take_profit"`
	Profit       float64     `json:"profit" db:"profit"`
	Commission   float64     `json:"commission" db:"commission"`
	Swap         float64     `json:"swap" db:"swap"`
	Status       TradeStatus `json:"status" db:"status"`
	OpenedAt     time.Time   `json:"opened_at" db:"opened_at"`
	ClosedAt     *time.Time  `json:"closed_at,omitempty" db:"closed_at"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
}

// CopiedTrade represents a trade copied to an investor's account
type CopiedTrade struct {
	ID               int64       `json:"id" db:"id"`
	OriginalTradeID  int64       `json:"original_trade_id" db:"original_trade_id"`
	SubscriptionUUID uuid.UUID   `json:"subscription_uuid" db:"subscription_uuid"`
	InvestorAccount  int64       `json:"investor_account_id" db:"investor_account_id"`
	Symbol           string      `json:"symbol" db:"symbol"`
	Type             TradeType   `json:"type" db:"type"`
	Volume           float64     `json:"volume" db:"volume"`
	CopyRatio        float64     `json:"copy_ratio" db:"copy_ratio"`
	OpenPrice        float64     `json:"open_price" db:"open_price"`
	ClosePrice       *float64    `json:"close_price,omitempty" db:"close_price"`
	StopLoss         *float64    `json:"stop_loss,omitempty" db:"stop_loss"`
	TakeProfit       *float64    `json:"take_profit,omitempty" db:"take_profit"`
	Profit           float64     `json:"profit" db:"profit"`
	Status           TradeStatus `json:"status" db:"status"`
	OpenedAt         time.Time   `json:"opened_at" db:"opened_at"`
	ClosedAt         *time.Time  `json:"closed_at,omitempty" db:"closed_at"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
}

// CreateTradeRequest represents the request to create a new trade
type CreateTradeRequest struct {
	StrategyUUID uuid.UUID `json:"strategy_uuid" binding:"required"`
	AccountID    int64     `json:"account_id" binding:"required"`
	Symbol       string    `json:"symbol" binding:"required"`
	Type         TradeType `json:"type" binding:"required,oneof=buy sell"`
	Volume       float64   `json:"volume" binding:"required,gt=0"`
	OpenPrice    float64   `json:"open_price" binding:"required,gt=0"`
	StopLoss     *float64  `json:"stop_loss,omitempty"`
	TakeProfit   *float64  `json:"take_profit,omitempty"`
}

// CopyTradeRequest represents the request to manually trigger trade copying
type CopyTradeRequest struct {
	SubscriptionUUIDs []uuid.UUID `json:"subscription_uuids,omitempty"`
}

// TradeFilter represents filter parameters for trade search
type TradeFilter struct {
	StrategyUUID uuid.UUID `form:"strategy_uuid"`
	common.TimeRange
	common.Pagination
}

// CopiedTradeFilter represents filter parameters for copied trade search
type CopiedTradeFilter struct {
	SubscriptionUUID uuid.UUID `form:"subscription_uuid"`
	common.Pagination
}

