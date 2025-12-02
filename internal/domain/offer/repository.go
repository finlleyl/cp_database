package offer

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Repository defines the interface for offer data operations
type Repository interface {
	Create(ctx context.Context, req *CreateOfferRequest) (*Offer, error)
	GetByID(ctx context.Context, id int64) (*Offer, error)
	List(ctx context.Context, filter *OfferFilter) (*common.PaginatedResult[Offer], error)
	Update(ctx context.Context, id int64, req *UpdateOfferRequest) (*Offer, error)
	ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest) (*Offer, error)
	GetByStrategyID(ctx context.Context, strategyID int64) ([]*Offer, error)
	GetActiveByStrategyID(ctx context.Context, strategyID int64) ([]*Offer, error)
}

type repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewRepository creates a new offer repository
func NewRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) Create(ctx context.Context, req *CreateOfferRequest) (*Offer, error) {
	query := `
		INSERT INTO offers (strategy_id, name, status, performance_fee_percent, management_fee_percent, registration_fee_amount)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, strategy_id, name, status, performance_fee_percent, management_fee_percent, registration_fee_amount, created_at, updated_at
	`

	var offer Offer
	err := r.db.QueryRowxContext(ctx, query,
		req.StrategyID,
		req.Name,
		common.OfferStatusActive,
		req.PerformanceFeePercent,
		req.ManagementFeePercent,
		req.RegistrationFeeAmount,
	).StructScan(&offer)
	if err != nil {
		r.logger.Error("Failed to create offer",
			zap.Int64("strategy_id", req.StrategyID),
			zap.String("name", req.Name),
			zap.Error(err))
		return nil, fmt.Errorf("create offer: %w", err)
	}

	r.logger.Info("Offer created",
		zap.Int64("id", offer.ID),
		zap.Int64("strategy_id", offer.StrategyID))

	return &offer, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*Offer, error) {
	query := `
		SELECT id, strategy_id, name, status, performance_fee_percent, management_fee_percent, registration_fee_amount, created_at, updated_at
		FROM offers
		WHERE id = $1
	`

	var offer Offer
	err := r.db.GetContext(ctx, &offer, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get offer by ID",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("get offer by id: %w", err)
	}

	return &offer, nil
}

func (r *repository) List(ctx context.Context, filter *OfferFilter) (*common.PaginatedResult[Offer], error) {
	filter.SetDefaults()

	var (
		conditions []string
		args       []interface{}
		argIndex   = 1
	)

	if filter.StrategyID != 0 {
		conditions = append(conditions, fmt.Sprintf("strategy_id = $%d", argIndex))
		args = append(args, filter.StrategyID)
		argIndex++
	}

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filter.Status)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM offers %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		r.logger.Error("Failed to count offers", zap.Error(err))
		return nil, fmt.Errorf("count offers: %w", err)
	}

	// Get data
	query := fmt.Sprintf(`
		SELECT id, strategy_id, name, status, performance_fee_percent, management_fee_percent, registration_fee_amount, created_at, updated_at
		FROM offers
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	var offers []Offer
	err = r.db.SelectContext(ctx, &offers, query, args...)
	if err != nil {
		r.logger.Error("Failed to list offers", zap.Error(err))
		return nil, fmt.Errorf("list offers: %w", err)
	}

	return &common.PaginatedResult[Offer]{
		Data:       offers,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(filter.Limit))),
	}, nil
}

func (r *repository) Update(ctx context.Context, id int64, req *UpdateOfferRequest) (*Offer, error) {
	var (
		setClauses []string
		args       []interface{}
		argIndex   = 1
	)

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}

	if req.PerformanceFeePercent != nil {
		setClauses = append(setClauses, fmt.Sprintf("performance_fee_percent = $%d", argIndex))
		args = append(args, *req.PerformanceFeePercent)
		argIndex++
	}

	if req.ManagementFeePercent != nil {
		setClauses = append(setClauses, fmt.Sprintf("management_fee_percent = $%d", argIndex))
		args = append(args, *req.ManagementFeePercent)
		argIndex++
	}

	if req.RegistrationFeeAmount != nil {
		setClauses = append(setClauses, fmt.Sprintf("registration_fee_amount = $%d", argIndex))
		args = append(args, *req.RegistrationFeeAmount)
		argIndex++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	setClauses = append(setClauses, "updated_at = now()")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE offers
		SET %s
		WHERE id = $%d
		RETURNING id, strategy_id, name, status, performance_fee_percent, management_fee_percent, registration_fee_amount, created_at, updated_at
	`, strings.Join(setClauses, ", "), argIndex)

	var offer Offer
	err := r.db.QueryRowxContext(ctx, query, args...).StructScan(&offer)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("offer not found: %d", id)
		}
		r.logger.Error("Failed to update offer",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("update offer: %w", err)
	}

	r.logger.Info("Offer updated", zap.Int64("id", offer.ID))

	return &offer, nil
}

func (r *repository) ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest) (*Offer, error) {
	query := `
		UPDATE offers
		SET status = $1, updated_at = now()
		WHERE id = $2
		RETURNING id, strategy_id, name, status, performance_fee_percent, management_fee_percent, registration_fee_amount, created_at, updated_at
	`

	var offer Offer
	err := r.db.QueryRowxContext(ctx, query, req.Status, id).StructScan(&offer)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("offer not found: %d", id)
		}
		r.logger.Error("Failed to change offer status",
			zap.Int64("id", id),
			zap.String("status", string(req.Status)),
			zap.Error(err))
		return nil, fmt.Errorf("change offer status: %w", err)
	}

	r.logger.Info("Offer status changed",
		zap.Int64("id", offer.ID),
		zap.String("status", string(offer.Status)))

	return &offer, nil
}

func (r *repository) GetByStrategyID(ctx context.Context, strategyID int64) ([]*Offer, error) {
	query := `
		SELECT id, strategy_id, name, status, performance_fee_percent, management_fee_percent, registration_fee_amount, created_at, updated_at
		FROM offers
		WHERE strategy_id = $1
		ORDER BY created_at DESC
	`

	var offers []*Offer
	err := r.db.SelectContext(ctx, &offers, query, strategyID)
	if err != nil {
		r.logger.Error("Failed to get offers by strategy ID",
			zap.Int64("strategy_id", strategyID),
			zap.Error(err))
		return nil, fmt.Errorf("get offers by strategy id: %w", err)
	}

	return offers, nil
}

func (r *repository) GetActiveByStrategyID(ctx context.Context, strategyID int64) ([]*Offer, error) {
	query := `
		SELECT id, strategy_id, name, status, performance_fee_percent, management_fee_percent, registration_fee_amount, created_at, updated_at
		FROM offers
		WHERE strategy_id = $1 AND status = 'active'
		ORDER BY created_at DESC
	`

	var offers []*Offer
	err := r.db.SelectContext(ctx, &offers, query, strategyID)
	if err != nil {
		r.logger.Error("Failed to get active offers by strategy ID",
			zap.Int64("strategy_id", strategyID),
			zap.Error(err))
		return nil, fmt.Errorf("get active offers by strategy id: %w", err)
	}

	return offers, nil
}
