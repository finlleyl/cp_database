package subscription

import (
	"context"
	"fmt"

	"github.com/finlleyl/cp_database/internal/domain/audit"
	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UseCase defines the interface for subscription business logic
type UseCase interface {
	Create(ctx context.Context, req *CreateSubscriptionRequest) (*Subscription, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*Subscription, error)
	List(ctx context.Context, filter *SubscriptionFilter) (*common.PaginatedResult[Subscription], error)
	Update(ctx context.Context, uuid uuid.UUID, req *UpdateSubscriptionRequest) (*Subscription, error)
	ChangeStatus(ctx context.Context, uuid uuid.UUID, req *ChangeStatusRequest, changedBy int64) (*Subscription, error)
	GetStatusHistory(ctx context.Context, uuid uuid.UUID) ([]*SubscriptionStatusHistory, error)
}

type useCase struct {
	repo      Repository
	auditRepo audit.Repository
	logger    *zap.Logger
}

// NewUseCase creates a new subscription use case
func NewUseCase(repo Repository, auditRepo audit.Repository, logger *zap.Logger) UseCase {
	return &useCase{repo: repo, auditRepo: auditRepo, logger: logger}
}

func (u *useCase) Create(ctx context.Context, req *CreateSubscriptionRequest) (*Subscription, error) {
	// TODO: Implement subscription creation business logic
	// Business flow:
	// 1. Validate offer exists and is active
	// 2. Validate investor account exists
	// 3. Validate risk_rules and config
	// 4. Create subscription with status = preparing
	// 5. Create status history record
	// 6. Create audit log
	u.logger.Info("UseCase: Creating subscription", 
		zap.Int64("investor_account_id", req.InvestorAccountID),
		zap.String("offer_uuid", req.OfferUUID.String()))
	
	subscription, err := u.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}
	
	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeSubscription,
		EntityID:   subscription.UUID.String(),
		Action:     audit.AuditActionCreate,
		NewValue:   subscription,
	})
	
	return subscription, nil
}

func (u *useCase) GetByUUID(ctx context.Context, uuid uuid.UUID) (*Subscription, error) {
	// TODO: Implement get subscription by UUID business logic
	u.logger.Info("UseCase: Getting subscription by UUID", zap.String("uuid", uuid.String()))
	return u.repo.GetByUUID(ctx, uuid)
}

func (u *useCase) List(ctx context.Context, filter *SubscriptionFilter) (*common.PaginatedResult[Subscription], error) {
	// TODO: Implement subscription listing business logic
	// Supports filtering by user_id and status
	filter.SetDefaults()
	u.logger.Info("UseCase: Listing subscriptions", zap.Any("filter", filter))
	return u.repo.List(ctx, filter)
}

func (u *useCase) Update(ctx context.Context, uuid uuid.UUID, req *UpdateSubscriptionRequest) (*Subscription, error) {
	// TODO: Implement subscription update business logic
	// 1. Get existing subscription
	// 2. Validate changes (config, risk_rules, filter)
	// 3. Update subscription
	// 4. Create audit log
	u.logger.Info("UseCase: Updating subscription", zap.String("uuid", uuid.String()))
	
	oldSubscription, err := u.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}
	
	subscription, err := u.repo.Update(ctx, uuid, req)
	if err != nil {
		return nil, fmt.Errorf("update subscription: %w", err)
	}
	
	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeSubscription,
		EntityID:   uuid.String(),
		Action:     audit.AuditActionUpdate,
		OldValue:   oldSubscription,
		NewValue:   subscription,
	})
	
	return subscription, nil
}

func (u *useCase) ChangeStatus(ctx context.Context, uuid uuid.UUID, req *ChangeStatusRequest, changedBy int64) (*Subscription, error) {
	// TODO: Implement subscription status change business logic
	// Business flow:
	// 1. Get existing subscription
	// 2. Validate status transition
	// 3. Update subscription status
	// 4. Create status history record
	// 5. Create audit log
	u.logger.Info("UseCase: Changing subscription status", 
		zap.String("uuid", uuid.String()),
		zap.String("status", string(req.Status)))
	
	oldSubscription, err := u.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}
	
	subscription, err := u.repo.ChangeStatus(ctx, uuid, req, changedBy)
	if err != nil {
		return nil, fmt.Errorf("change subscription status: %w", err)
	}
	
	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeSubscription,
		EntityID:   uuid.String(),
		Action:     audit.AuditActionStatusChange,
		OldValue:   oldSubscription,
		NewValue:   subscription,
		Changes: map[string]interface{}{
			"old_status":    oldSubscription.Status,
			"new_status":    subscription.Status,
			"status_reason": req.StatusReason,
		},
	})
	
	return subscription, nil
}

func (u *useCase) GetStatusHistory(ctx context.Context, uuid uuid.UUID) ([]*SubscriptionStatusHistory, error) {
	// TODO: Implement get subscription status history business logic
	u.logger.Info("UseCase: Getting subscription status history", zap.String("uuid", uuid.String()))
	return u.repo.GetStatusHistory(ctx, uuid)
}

