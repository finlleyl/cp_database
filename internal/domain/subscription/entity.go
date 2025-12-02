package subscription

import (
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

// Subscription represents an investor's subscription to a strategy
type Subscription struct {
	ID                int64                     `json:"id" db:"id"`
	InvestorUserID    int64                     `json:"investor_user_id" db:"investor_user_id"`
	InvestorAccountID int64                     `json:"investor_account_id" db:"investor_account_id"`
	OfferID           int64                     `json:"offer_id" db:"offer_id"`
	Status            common.SubscriptionStatus `json:"status" db:"status"`
	CreatedAt         time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at" db:"updated_at"`
}

// CreateSubscriptionRequest represents the request to create a new subscription
type CreateSubscriptionRequest struct {
	InvestorUserID    int64 `json:"investor_user_id" binding:"required"`
	InvestorAccountID int64 `json:"investor_account_id" binding:"required"`
	OfferID           int64 `json:"offer_id" binding:"required"`
}

// UpdateSubscriptionRequest represents the request to update a subscription
type UpdateSubscriptionRequest struct {
	// Currently no updatable fields besides status
}

// ChangeStatusRequest represents the request to change subscription status
type ChangeStatusRequest struct {
	Status       common.SubscriptionStatus `json:"status" binding:"required,oneof=active archived suspended deleted"`
	StatusReason string                    `json:"status_reason"`
}

// SubscriptionStatusHistory represents a status change history record
type SubscriptionStatusHistory struct {
	ID             int64                     `json:"id" db:"id"`
	SubscriptionID int64                     `json:"subscription_id" db:"subscription_id"`
	OldStatus      common.SubscriptionStatus `json:"old_status" db:"old_status"`
	NewStatus      common.SubscriptionStatus `json:"new_status" db:"new_status"`
	Reason         string                    `json:"reason" db:"reason"`
	ChangedBy      int64                     `json:"changed_by" db:"changed_by"`
	CreatedAt      time.Time                 `json:"created_at" db:"created_at"`
}

// SubscriptionFilter represents filter parameters for subscription search
type SubscriptionFilter struct {
	UserID  int64                     `form:"user_id"`
	OfferID int64                     `form:"offer_id"`
	Status  common.SubscriptionStatus `form:"status"`
	common.Pagination
}
