package offer

import (
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

type Offer struct {
	ID                    int64              `json:"id" db:"id"`
	StrategyID            int64              `json:"strategy_id" db:"strategy_id"`
	Name                  string             `json:"name" db:"name"`
	Status                common.OfferStatus `json:"status" db:"status"`
	PerformanceFeePercent *float64           `json:"performance_fee_percent" db:"performance_fee_percent"`
	ManagementFeePercent  *float64           `json:"management_fee_percent" db:"management_fee_percent"`
	RegistrationFeeAmount *float64           `json:"registration_fee_amount" db:"registration_fee_amount"`
	CreatedAt             time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at" db:"updated_at"`
}

type CreateOfferRequest struct {
	StrategyID            int64    `json:"strategy_id" binding:"required"`
	Name                  string   `json:"name" binding:"required"`
	PerformanceFeePercent *float64 `json:"performance_fee_percent"`
	ManagementFeePercent  *float64 `json:"management_fee_percent"`
	RegistrationFeeAmount *float64 `json:"registration_fee_amount"`
}

type UpdateOfferRequest struct {
	Name                  *string  `json:"name,omitempty"`
	PerformanceFeePercent *float64 `json:"performance_fee_percent,omitempty"`
	ManagementFeePercent  *float64 `json:"management_fee_percent,omitempty"`
	RegistrationFeeAmount *float64 `json:"registration_fee_amount,omitempty"`
}

type ChangeStatusRequest struct {
	Status common.OfferStatus `json:"status" binding:"required,oneof=active archived deleted"`
}

type OfferFilter struct {
	StrategyID int64              `form:"strategy_id"`
	Status     common.OfferStatus `form:"status"`
	common.Pagination
}
