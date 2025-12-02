package offer

import (
	"context"
	"fmt"

	"github.com/finlleyl/cp_database/internal/domain/audit"
	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/finlleyl/cp_database/internal/domain/strategy"
	"go.uber.org/zap"
)

// UseCase defines the interface for offer business logic
type UseCase interface {
	Create(ctx context.Context, req *CreateOfferRequest) (*Offer, error)
	GetByID(ctx context.Context, id int64) (*Offer, error)
	List(ctx context.Context, filter *OfferFilter) (*common.PaginatedResult[Offer], error)
	Update(ctx context.Context, id int64, req *UpdateOfferRequest) (*Offer, error)
	ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest) (*Offer, error)
}

type useCase struct {
	repo         Repository
	strategyRepo strategy.Repository
	auditRepo    audit.Repository
	logger       *zap.Logger
}

// NewUseCase creates a new offer use case
func NewUseCase(
	repo Repository,
	strategyRepo strategy.Repository,
	auditRepo audit.Repository,
	logger *zap.Logger,
) UseCase {
	return &useCase{
		repo:         repo,
		strategyRepo: strategyRepo,
		auditRepo:    auditRepo,
		logger:       logger,
	}
}

func (u *useCase) Create(ctx context.Context, req *CreateOfferRequest) (*Offer, error) {
	u.logger.Info("UseCase: Creating offer",
		zap.Int64("strategy_id", req.StrategyID),
		zap.String("name", req.Name))

	// Validate strategy exists
	strat, err := u.strategyRepo.GetByID(ctx, req.StrategyID)
	if err != nil {
		return nil, fmt.Errorf("get strategy: %w", err)
	}
	if strat == nil {
		return nil, fmt.Errorf("strategy not found")
	}

	offer, err := u.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create offer: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeOffer,
		EntityID:   offer.ID,
		Action:     audit.AuditActionCreate,
		NewValue:   offer,
	})

	return offer, nil
}

func (u *useCase) GetByID(ctx context.Context, id int64) (*Offer, error) {
	u.logger.Info("UseCase: Getting offer by ID", zap.Int64("id", id))
	return u.repo.GetByID(ctx, id)
}

func (u *useCase) List(ctx context.Context, filter *OfferFilter) (*common.PaginatedResult[Offer], error) {
	filter.SetDefaults()
	u.logger.Info("UseCase: Listing offers", zap.Any("filter", filter))
	return u.repo.List(ctx, filter)
}

func (u *useCase) Update(ctx context.Context, id int64, req *UpdateOfferRequest) (*Offer, error) {
	u.logger.Info("UseCase: Updating offer", zap.Int64("id", id))

	oldOffer, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get offer: %w", err)
	}

	offer, err := u.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("update offer: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeOffer,
		EntityID:   offer.ID,
		Action:     audit.AuditActionUpdate,
		OldValue:   oldOffer,
		NewValue:   offer,
	})

	return offer, nil
}

func (u *useCase) ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest) (*Offer, error) {
	u.logger.Info("UseCase: Changing offer status",
		zap.Int64("id", id),
		zap.String("status", string(req.Status)))

	oldOffer, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get offer: %w", err)
	}

	offer, err := u.repo.ChangeStatus(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("change offer status: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeOffer,
		EntityID:   offer.ID,
		Action:     audit.AuditActionStatusChange,
		OldValue:   oldOffer,
		NewValue:   offer,
		Changes: map[string]interface{}{
			"old_status": oldOffer.Status,
			"new_status": offer.Status,
		},
	})

	return offer, nil
}
