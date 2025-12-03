package audit

import (
	"context"
	"fmt"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"go.uber.org/zap"
)

type UseCase interface {

	List(ctx context.Context, filter *AuditFilter) (*common.PaginatedResult[AuditLog], error)

	GetByEntity(ctx context.Context, entityName string, entityPK string) ([]*AuditLog, error)

	GetStats(ctx context.Context, filter *AuditStatsFilter) ([]*AuditStats, error)
}

type useCase struct {
	repo   Repository
	logger *zap.Logger
}

func NewUseCase(repo Repository, logger *zap.Logger) UseCase {
	return &useCase{repo: repo, logger: logger}
}

func (u *useCase) List(ctx context.Context, filter *AuditFilter) (*common.PaginatedResult[AuditLog], error) {
	filter.SetDefaults()

	u.logger.Debug("Listing audit logs",
		zap.String("entity_name", filter.EntityName),
		zap.String("entity_pk", filter.EntityPK),
		zap.String("operation", string(filter.Operation)))

	logs, err := u.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}

	return logs, nil
}

func (u *useCase) GetByEntity(ctx context.Context, entityName string, entityPK string) ([]*AuditLog, error) {

	validEntities := map[string]bool{
		EntityNameUsers:         true,
		EntityNameAccounts:      true,
		EntityNameStrategies:    true,
		EntityNameOffers:        true,
		EntityNameSubscriptions: true,
		EntityNameTrades:        true,
	}

	if !validEntities[entityName] {
		return nil, fmt.Errorf("invalid entity name: %s", entityName)
	}

	u.logger.Debug("Getting audit logs by entity",
		zap.String("entity_name", entityName),
		zap.String("entity_pk", entityPK))

	logs, err := u.repo.GetByEntity(ctx, entityName, entityPK)
	if err != nil {
		return nil, fmt.Errorf("get audit logs by entity: %w", err)
	}

	return logs, nil
}

func (u *useCase) GetStats(ctx context.Context, filter *AuditStatsFilter) ([]*AuditStats, error) {
	u.logger.Debug("Getting audit stats",
		zap.String("entity_name", filter.EntityName),
		zap.Time("from", filter.From),
		zap.Time("to", filter.To))

	stats, err := u.repo.GetStats(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get audit stats: %w", err)
	}

	return stats, nil
}
