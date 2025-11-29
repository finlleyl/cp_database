package offer

import (
	"context"
	"fmt"

	"github.com/finlleyl/cp_database/internal/domain/audit"
	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/finlleyl/cp_database/internal/domain/strategy"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UseCase defines the interface for offer business logic
type UseCase interface {
	Create(ctx context.Context, req *CreateOfferRequest) (*Offer, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*Offer, error)
	List(ctx context.Context, filter *OfferFilter) (*common.PaginatedResult[Offer], error)
	Update(ctx context.Context, uuid uuid.UUID, req *UpdateOfferRequest) (*Offer, error)
	ChangeStatus(ctx context.Context, uuid uuid.UUID, req *ChangeStatusRequest) (*Offer, error)
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
	// TODO: Implement offer creation business logic
	// Business flow:
	// 1. Validate strategy exists and is active
	// 2. Validate fee settings
	// 3. Create offer with status = active
	// 4. Create audit log
	u.logger.Info("UseCase: Creating offer", 
		zap.String("strategy_uuid", req.StrategyUUID.String()),
		zap.String("name", req.Name))
	
	// Validate strategy exists
	strat, err := u.strategyRepo.GetByUUID(ctx, req.StrategyUUID)
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
		EntityID:   offer.UUID.String(),
		Action:     audit.AuditActionCreate,
		NewValue:   offer,
	})
	
	return offer, nil
}

func (u *useCase) GetByUUID(ctx context.Context, uuid uuid.UUID) (*Offer, error) {
	// TODO: Implement get offer by UUID business logic
	u.logger.Info("UseCase: Getting offer by UUID", zap.String("uuid", uuid.String()))
	return u.repo.GetByUUID(ctx, uuid)
}

func (u *useCase) List(ctx context.Context, filter *OfferFilter) (*common.PaginatedResult[Offer], error) {
	// TODO: Implement offer listing business logic
	// Supports filtering by strategy_uuid and status
	filter.SetDefaults()
	u.logger.Info("UseCase: Listing offers", zap.Any("filter", filter))
	return u.repo.List(ctx, filter)
}

func (u *useCase) Update(ctx context.Context, uuid uuid.UUID, req *UpdateOfferRequest) (*Offer, error) {
	// TODO: Implement offer update business logic
	// 1. Get existing offer
	// 2. Validate changes (may restrict changes if active subscriptions exist)
	// 3. Update offer
	// 4. Create audit log
	u.logger.Info("UseCase: Updating offer", zap.String("uuid", uuid.String()))
	
	oldOffer, err := u.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("get offer: %w", err)
	}
	
	offer, err := u.repo.Update(ctx, uuid, req)
	if err != nil {
		return nil, fmt.Errorf("update offer: %w", err)
	}
	
	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeOffer,
		EntityID:   uuid.String(),
		Action:     audit.AuditActionUpdate,
		OldValue:   oldOffer,
		NewValue:   offer,
	})
	
	return offer, nil
}

func (u *useCase) ChangeStatus(ctx context.Context, uuid uuid.UUID, req *ChangeStatusRequest) (*Offer, error) {
	// TODO: Implement offer status change business logic
	u.logger.Info("UseCase: Changing offer status", 
		zap.String("uuid", uuid.String()),
		zap.String("status", string(req.Status)))
	
	oldOffer, err := u.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("get offer: %w", err)
	}
	
	offer, err := u.repo.ChangeStatus(ctx, uuid, req)
	if err != nil {
		return nil, fmt.Errorf("change offer status: %w", err)
	}
	
	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeOffer,
		EntityID:   uuid.String(),
		Action:     audit.AuditActionStatusChange,
		OldValue:   oldOffer,
		NewValue:   offer,
		Changes: map[string]interface{}{
			"old_status":    oldOffer.Status,
			"new_status":    offer.Status,
			"status_reason": req.StatusReason,
		},
	})
	
	return offer, nil
}

