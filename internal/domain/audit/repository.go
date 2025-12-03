package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Repository interface {

	Create(ctx context.Context, req *AuditCreateRequest) (*AuditLog, error)

	List(ctx context.Context, filter *AuditFilter) (*common.PaginatedResult[AuditLog], error)

	GetByEntity(ctx context.Context, entityName string, entityPK string) ([]*AuditLog, error)

	GetStats(ctx context.Context, filter *AuditStatsFilter) ([]*AuditStats, error)

	CountByEntity(ctx context.Context, entityName string) (int64, error)
}

type repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) Create(ctx context.Context, req *AuditCreateRequest) (*AuditLog, error) {
	var oldRowJSON, newRowJSON []byte
	var err error

	if req.OldValue != nil {
		oldRowJSON, err = json.Marshal(req.OldValue)
		if err != nil {
			r.logger.Warn("Failed to marshal old value for audit",
				zap.Error(err))
			oldRowJSON = nil
		}
	}

	if req.NewValue != nil {
		newRowJSON, err = json.Marshal(req.NewValue)
		if err != nil {
			r.logger.Warn("Failed to marshal new value for audit",
				zap.Error(err))
			newRowJSON = nil
		}
	}

	query := `
		INSERT INTO audit_log (entity_name, entity_pk, operation, changed_by, old_row, new_row)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, entity_name, entity_pk, operation, changed_by, changed_at, old_row, new_row
	`

	var result AuditLog
	err = r.db.QueryRowxContext(ctx, query,
		string(req.EntityType),
		fmt.Sprintf("%d", req.EntityID),
		string(req.Action),
		req.UserID,
		oldRowJSON,
		newRowJSON,
	).StructScan(&result)

	if err != nil {
		r.logger.Error("Failed to create audit log",
			zap.String("entity_type", string(req.EntityType)),
			zap.Int64("entity_id", req.EntityID),
			zap.Error(err))
		return nil, fmt.Errorf("create audit log: %w", err)
	}

	r.logger.Debug("Audit log created",
		zap.Int64("id", result.ID),
		zap.String("entity_name", result.EntityName),
		zap.String("entity_pk", result.EntityPK))

	return &result, nil
}

func (r *repository) List(ctx context.Context, filter *AuditFilter) (*common.PaginatedResult[AuditLog], error) {
	filter.SetDefaults()

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.EntityName != "" {
		conditions = append(conditions, fmt.Sprintf("entity_name = $%d", argIndex))
		args = append(args, filter.EntityName)
		argIndex++
	}

	if filter.EntityPK != "" {
		conditions = append(conditions, fmt.Sprintf("entity_pk = $%d", argIndex))
		args = append(args, filter.EntityPK)
		argIndex++
	}

	if filter.Operation != "" {
		conditions = append(conditions, fmt.Sprintf("operation = $%d", argIndex))
		args = append(args, filter.Operation)
		argIndex++
	}

	if filter.ChangedBy != nil {
		conditions = append(conditions, fmt.Sprintf("changed_by = $%d", argIndex))
		args = append(args, *filter.ChangedBy)
		argIndex++
	}

	if !filter.From.IsZero() {
		conditions = append(conditions, fmt.Sprintf("changed_at >= $%d", argIndex))
		args = append(args, filter.From)
		argIndex++
	}

	if !filter.To.IsZero() {
		conditions = append(conditions, fmt.Sprintf("changed_at <= $%d", argIndex))
		args = append(args, filter.To)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_log %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		r.logger.Error("Failed to count audit logs", zap.Error(err))
		return nil, fmt.Errorf("count audit logs: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, entity_name, entity_pk, operation, changed_by, changed_at, old_row, new_row
		FROM audit_log
		%s
		ORDER BY changed_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	var logs []AuditLog
	err = r.db.SelectContext(ctx, &logs, query, args...)
	if err != nil {
		r.logger.Error("Failed to list audit logs", zap.Error(err))
		return nil, fmt.Errorf("list audit logs: %w", err)
	}

	totalPages := int(total) / filter.Limit
	if int(total)%filter.Limit > 0 {
		totalPages++
	}

	return &common.PaginatedResult[AuditLog]{
		Data:       logs,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

func (r *repository) GetByEntity(ctx context.Context, entityName string, entityPK string) ([]*AuditLog, error) {
	query := `
		SELECT id, entity_name, entity_pk, operation, changed_by, changed_at, old_row, new_row
		FROM audit_log
		WHERE entity_name = $1 AND entity_pk = $2
		ORDER BY changed_at DESC
	`

	var logs []*AuditLog
	err := r.db.SelectContext(ctx, &logs, query, entityName, entityPK)
	if err != nil {
		r.logger.Error("Failed to get audit logs by entity",
			zap.String("entity_name", entityName),
			zap.String("entity_pk", entityPK),
			zap.Error(err))
		return nil, fmt.Errorf("get audit logs by entity: %w", err)
	}

	return logs, nil
}

func (r *repository) GetStats(ctx context.Context, filter *AuditStatsFilter) ([]*AuditStats, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.EntityName != "" {
		conditions = append(conditions, fmt.Sprintf("entity_name = $%d", argIndex))
		args = append(args, filter.EntityName)
		argIndex++
	}

	if !filter.From.IsZero() {
		conditions = append(conditions, fmt.Sprintf("changed_at >= $%d", argIndex))
		args = append(args, filter.From)
		argIndex++
	}

	if !filter.To.IsZero() {
		conditions = append(conditions, fmt.Sprintf("changed_at <= $%d", argIndex))
		args = append(args, filter.To)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT entity_name, operation::text, COUNT(*) as total_changes
		FROM audit_log
		%s
		GROUP BY entity_name, operation
		ORDER BY entity_name, operation
	`, whereClause)

	var stats []*AuditStats
	err := r.db.SelectContext(ctx, &stats, query, args...)
	if err != nil {
		r.logger.Error("Failed to get audit stats", zap.Error(err))
		return nil, fmt.Errorf("get audit stats: %w", err)
	}

	return stats, nil
}

func (r *repository) CountByEntity(ctx context.Context, entityName string) (int64, error) {
	query := `SELECT COUNT(*) FROM audit_log WHERE entity_name = $1`

	var count int64
	err := r.db.GetContext(ctx, &count, query, entityName)
	if err != nil {
		r.logger.Error("Failed to count audit logs by entity",
			zap.String("entity_name", entityName),
			zap.Error(err))
		return 0, fmt.Errorf("count audit logs by entity: %w", err)
	}

	return count, nil
}
