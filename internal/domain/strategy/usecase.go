package strategy

import (
	"context"
	"fmt"

	"github.com/finlleyl/cp_database/internal/domain/audit"
	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/finlleyl/cp_database/internal/domain/subscription"
	"go.uber.org/zap"
)

// UseCase defines the interface for strategy business logic
type UseCase interface {
	Create(ctx context.Context, req *CreateStrategyRequest) (*Strategy, error)
	GetByID(ctx context.Context, id int64) (*GetStrategyByIDResponse, error)
	List(ctx context.Context, filter *StrategyFilter) (*common.PaginatedResult[Strategy], error)
	Update(ctx context.Context, id int64, req *UpdateStrategyRequest) (*Strategy, error)
	ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest) (*Strategy, error)
}

type useCase struct {
	repo             Repository
	subscriptionRepo subscription.Repository
	auditRepo        audit.Repository
	logger           *zap.Logger
}

// NewUseCase creates a new strategy use case
func NewUseCase(
	repo Repository,
	subscriptionRepo subscription.Repository,
	auditRepo audit.Repository,
	logger *zap.Logger,
) UseCase {
	return &useCase{
		repo:             repo,
		subscriptionRepo: subscriptionRepo,
		auditRepo:        auditRepo,
		logger:           logger,
	}
}

func (u *useCase) Create(ctx context.Context, req *CreateStrategyRequest) (*Strategy, error) {
	// TODO: Implement strategy creation business logic
	// Business flow:
	// 1. Validate account exists and belongs to user
	// 2. Check if account doesn't have existing active strategy
	// 3. Create strategy with status = preparing
	// 4. Create audit log
	u.logger.Info("UseCase: Creating strategy",
		zap.String("nickname", req.Nickname),
		zap.Int64("account_id", req.AccountID))

	strategy, err := u.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create strategy: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeStrategy,
		EntityID:   strategy.ID,
		Action:     audit.AuditActionCreate,
		NewValue:   strategy,
	})

	return strategy, nil
}

func (u *useCase) GetByID(ctx context.Context, id int64) (*GetStrategyByIDResponse, error) {
	// TODO: Implement get strategy by ID business logic
	u.logger.Info("UseCase: Getting strategy by ID", zap.Int64("id", id))
	return u.repo.GetByID(ctx, id)
}

func (u *useCase) List(ctx context.Context, filter *StrategyFilter) (*common.PaginatedResult[Strategy], error) {
	// TODO: Implement strategy listing business logic
	// Supports filtering by status, min_roi, max_drawdown_pct, risk_score, search by nickname
	filter.SetDefaults()
	u.logger.Info("UseCase: Listing strategies", zap.Any("filter", filter))
	return u.repo.List(ctx, filter)
}

func (u *useCase) Update(ctx context.Context, id int64, req *UpdateStrategyRequest) (*Strategy, error) {
	// TODO: Implement strategy update business logic
	// 1. Get existing strategy
	// 2. Validate changes
	// 3. Update strategy
	// 4. Create audit log
	u.logger.Info("UseCase: Updating strategy", zap.Int64("id", id))

	oldStrategy, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get strategy: %w", err)
	}

	strategy, err := u.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("update strategy: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeStrategy,
		EntityID:   id,
		Action:     audit.AuditActionUpdate,
		OldValue:   oldStrategy,
		NewValue:   strategy,
	})

	return strategy, nil
}

func (u *useCase) ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest) (*Strategy, error) {
	// TODO: Implement strategy status change business logic
	// Business flow for archived/deleted:
	// 1. Get existing strategy
	// 2. Validate status transition
	// 3. If archiving/deleting: archive all active subscriptions with reason
	// 4. Update strategy status
	// 5. Create audit log
	u.logger.Info("UseCase: Changing strategy status",
		zap.Int64("id", id),
		zap.String("status", string(req.Status)))

	oldStrategy, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get strategy: %w", err)
	}

	// If archiving or deleting, archive all active subscriptions
	if req.Status == common.StrategyStatusArchived || req.Status == common.StrategyStatusDeleted {
		reason := fmt.Sprintf("strategy_%s", req.Status)
		if req.StatusReason != "" {
			reason = req.StatusReason
		}
		if err := u.subscriptionRepo.ArchiveByStrategyID(ctx, id, reason); err != nil {
			u.logger.Warn("Failed to archive subscriptions", zap.Error(err))
		}
	}

	strategy, err := u.repo.ChangeStatus(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("change strategy status: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeStrategy,
		EntityID:   strategy.ID,
		Action:     audit.AuditActionStatusChange,
		OldValue:   oldStrategy,
		NewValue:   strategy,
		Changes: map[string]interface{}{
			"old_status":    oldStrategy.Status,
			"new_status":    strategy.Status,
			"status_reason": req.StatusReason,
		},
	})

	return strategy, nil
}
