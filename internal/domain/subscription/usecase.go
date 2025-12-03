package subscription

import (
	"context"
	"fmt"

	"github.com/finlleyl/cp_database/internal/domain/audit"
	"github.com/finlleyl/cp_database/internal/domain/common"
	"go.uber.org/zap"
)

type UseCase interface {
	Create(ctx context.Context, req *CreateSubscriptionRequest) (*Subscription, error)
	GetByID(ctx context.Context, id int64) (*Subscription, error)
	List(ctx context.Context, filter *SubscriptionFilter) (*common.PaginatedResult[Subscription], error)
	Update(ctx context.Context, id int64, req *UpdateSubscriptionRequest) (*Subscription, error)
	ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest, changedBy int64) (*Subscription, error)
	GetStatusHistory(ctx context.Context, id int64) ([]*SubscriptionStatusHistory, error)
}

type useCase struct {
	repo      Repository
	auditRepo audit.Repository
	logger    *zap.Logger
}

func NewUseCase(repo Repository, auditRepo audit.Repository, logger *zap.Logger) UseCase {
	return &useCase{repo: repo, auditRepo: auditRepo, logger: logger}
}

func (u *useCase) Create(ctx context.Context, req *CreateSubscriptionRequest) (*Subscription, error) {

	u.logger.Info("UseCase: Creating subscription",
		zap.Int64("investor_account_id", req.InvestorAccountID),
		zap.Int64("offer_id", req.OfferID))

	subscription, err := u.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeSubscription,
		EntityID:   subscription.ID,
		Action:     audit.AuditActionCreate,
		NewValue:   subscription,
	})

	return subscription, nil
}

func (u *useCase) GetByID(ctx context.Context, id int64) (*Subscription, error) {

	u.logger.Info("UseCase: Getting subscription by ID", zap.Int64("id", id))
	return u.repo.GetByID(ctx, id)
}

func (u *useCase) List(ctx context.Context, filter *SubscriptionFilter) (*common.PaginatedResult[Subscription], error) {

	filter.SetDefaults()
	u.logger.Info("UseCase: Listing subscriptions", zap.Any("filter", filter))
	return u.repo.List(ctx, filter)
}

func (u *useCase) Update(ctx context.Context, id int64, req *UpdateSubscriptionRequest) (*Subscription, error) {

	u.logger.Info("UseCase: Updating subscription", zap.Int64("id", id))

	oldSubscription, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}

	subscription, err := u.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("update subscription: %w", err)
	}

	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeSubscription,
		EntityID:   subscription.ID,
		Action:     audit.AuditActionUpdate,
		OldValue:   oldSubscription,
		NewValue:   subscription,
	})

	return subscription, nil
}

func (u *useCase) ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest, changedBy int64) (*Subscription, error) {

	u.logger.Info("UseCase: Changing subscription status",
		zap.Int64("id", id),
		zap.String("status", string(req.Status)))

	oldSubscription, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}

	subscription, err := u.repo.ChangeStatus(ctx, id, req, changedBy)
	if err != nil {
		return nil, fmt.Errorf("change subscription status: %w", err)
	}

	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeSubscription,
		EntityID:   subscription.ID,
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

func (u *useCase) GetStatusHistory(ctx context.Context, id int64) ([]*SubscriptionStatusHistory, error) {

	u.logger.Info("UseCase: Getting subscription status history", zap.Int64("id", id))
	return u.repo.GetStatusHistory(ctx, id)
}
