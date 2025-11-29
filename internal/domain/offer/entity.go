package offer

import (
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/google/uuid"
)

// Offer represents a strategy's offer for investors
type Offer struct {
	UUID                   uuid.UUID          `json:"uuid" db:"uuid"`
	StrategyUUID           uuid.UUID          `json:"strategy_uuid" db:"strategy_uuid"`
	Name                   string             `json:"name" db:"name"`
	PerformanceFee         float64            `json:"performance_fee" db:"performance_fee"`
	PerformanceFeeInterval common.FeeInterval `json:"performance_fee_interval" db:"performance_fee_interval"`
	ManagementFee          float64            `json:"management_fee" db:"management_fee"`
	ManagementFeeInterval  common.FeeInterval `json:"management_fee_interval" db:"management_fee_interval"`
	RegistrationFee        float64            `json:"registration_fee" db:"registration_fee"`
	Status                 common.OfferStatus `json:"status" db:"status"`
	StatusReason           string             `json:"status_reason" db:"status_reason"`
	CreatedAt              time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at" db:"updated_at"`
}

// CreateOfferRequest represents the request to create a new offer
type CreateOfferRequest struct {
	StrategyUUID           uuid.UUID          `json:"strategy_uuid" binding:"required"`
	Name                   string             `json:"name" binding:"required"`
	PerformanceFee         float64            `json:"performance_fee" binding:"gte=0,lte=100"`
	PerformanceFeeInterval common.FeeInterval `json:"performance_fee_interval" binding:"required,oneof=daily weekly monthly"`
	ManagementFee          float64            `json:"management_fee" binding:"gte=0,lte=100"`
	ManagementFeeInterval  common.FeeInterval `json:"management_fee_interval" binding:"required,oneof=daily weekly monthly"`
	RegistrationFee        float64            `json:"registration_fee" binding:"gte=0"`
}

// UpdateOfferRequest represents the request to update an offer
type UpdateOfferRequest struct {
	Name                   *string             `json:"name,omitempty"`
	PerformanceFee         *float64            `json:"performance_fee,omitempty"`
	PerformanceFeeInterval *common.FeeInterval `json:"performance_fee_interval,omitempty"`
	ManagementFee          *float64            `json:"management_fee,omitempty"`
	ManagementFeeInterval  *common.FeeInterval `json:"management_fee_interval,omitempty"`
	RegistrationFee        *float64            `json:"registration_fee,omitempty"`
}

// ChangeStatusRequest represents the request to change offer status
type ChangeStatusRequest struct {
	Status       common.OfferStatus `json:"status" binding:"required,oneof=active archived deleted"`
	StatusReason string             `json:"status_reason"`
}

// OfferFilter represents filter parameters for offer search
type OfferFilter struct {
	StrategyUUID uuid.UUID          `form:"strategy_uuid"`
	Status       common.OfferStatus `form:"status"`
	common.Pagination
}
