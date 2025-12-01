package audit

import (
	"context"
	"fmt"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"go.uber.org/zap"
)

// UseCase defines the interface for audit log business logic
type UseCase interface {
	List(ctx context.Context, filter *AuditFilter) (*common.PaginatedResult[AuditLog], error)
	GetByEntityID(ctx context.Context, entityType EntityType, entityID string) ([]*AuditLog, error)
}

type useCase struct {
	repo   Repository
	logger *zap.Logger
}

// NewUseCase creates a new audit use case
func NewUseCase(repo Repository, logger *zap.Logger) UseCase {
	return &useCase{repo: repo, logger: logger}
}

func (u *useCase) List(ctx context.Context, filter *AuditFilter) (*common.PaginatedResult[AuditLog], error) {
	// TODO: Implement audit log listing business logic
	// Supports filtering by entity, entity_id, action, user_id, time range
	filter.SetDefaults()
	u.logger.Info("UseCase: Listing audit logs", zap.Any("filter", filter))

	logs, err := u.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}

	return logs, nil
}

func (u *useCase) GetByEntityID(ctx context.Context, entityType EntityType, entityID string) ([]*AuditLog, error) {
	// TODO: Implement get audit logs by entity business logic
	u.logger.Info("UseCase: Getting audit logs by entity",
		zap.String("entity_type", string(entityType)),
		zap.String("entity_id", entityID))

	logs, err := u.repo.GetByEntityID(ctx, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("get audit logs by entity: %w", err)
	}

	return logs, nil
}
