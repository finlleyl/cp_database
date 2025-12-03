package trade

import (
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

type TradeDirection string

const (
	TradeDirectionBuy  TradeDirection = "buy"
	TradeDirectionSell TradeDirection = "sell"
)

type Trade struct {
	ID              int64          `json:"id" db:"id"`
	StrategyID      int64          `json:"strategy_id" db:"strategy_id"`
	MasterAccountID int64          `json:"master_account_id" db:"master_account_id"`
	Symbol          string         `json:"symbol" db:"symbol"`
	VolumeLots      float64        `json:"volume_lots" db:"volume_lots"`
	Direction       TradeDirection `json:"direction" db:"direction"`
	OpenTime        time.Time      `json:"open_time" db:"open_time"`
	CloseTime       *time.Time     `json:"close_time,omitempty" db:"close_time"`
	OpenPrice       float64        `json:"open_price" db:"open_price"`
	ClosePrice      *float64       `json:"close_price,omitempty" db:"close_price"`
	Profit          *float64       `json:"profit,omitempty" db:"profit"`
	Commission      *float64       `json:"commission,omitempty" db:"commission"`
	Swap            *float64       `json:"swap,omitempty" db:"swap"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
}

type CopiedTrade struct {
	ID                int64      `json:"id" db:"id"`
	TradeID           int64      `json:"trade_id" db:"trade_id"`
	SubscriptionID    int64      `json:"subscription_id" db:"subscription_id"`
	InvestorAccountID int64      `json:"investor_account_id" db:"investor_account_id"`
	VolumeLots        float64    `json:"volume_lots" db:"volume_lots"`
	Profit            *float64   `json:"profit,omitempty" db:"profit"`
	Commission        *float64   `json:"commission,omitempty" db:"commission"`
	Swap              *float64   `json:"swap,omitempty" db:"swap"`
	OpenTime          time.Time  `json:"open_time" db:"open_time"`
	CloseTime         *time.Time `json:"close_time,omitempty" db:"close_time"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}

type CreateTradeRequest struct {
	StrategyID      int64          `json:"strategy_id" binding:"required"`
	MasterAccountID int64          `json:"master_account_id" binding:"required"`
	Symbol          string         `json:"symbol" binding:"required"`
	VolumeLots      float64        `json:"volume_lots" binding:"required,gt=0"`
	Direction       TradeDirection `json:"direction" binding:"required,oneof=buy sell"`
	OpenTime        time.Time      `json:"open_time" binding:"required"`
	OpenPrice       float64        `json:"open_price" binding:"required,gt=0"`
}

type CreateCopiedTradeRequest struct {
	TradeID           int64     `json:"trade_id" binding:"required"`
	SubscriptionID    int64     `json:"subscription_id" binding:"required"`
	InvestorAccountID int64     `json:"investor_account_id" binding:"required"`
	VolumeLots        float64   `json:"volume_lots" binding:"required,gt=0"`
	OpenTime          time.Time `json:"open_time" binding:"required"`
}

type TradeFilter struct {
	StrategyID int64 `form:"strategy_id"`
	common.TimeRange
	common.Pagination
}

type CopiedTradeFilter struct {
	SubscriptionID int64 `form:"subscription_id"`
	TradeID        int64 `form:"trade_id"`
	common.Pagination
}

type CopyTradeRequest struct {
	SubscriptionIDs []int64 `json:"subscription_ids,omitempty"`
}
