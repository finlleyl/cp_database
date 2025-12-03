package subscription

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

type Repository interface {
	Create(ctx context.Context, req *CreateSubscriptionRequest) (*Subscription, error)
	GetByID(ctx context.Context, id int64) (*Subscription, error)
	List(ctx context.Context, filter *SubscriptionFilter) (*common.PaginatedResult[Subscription], error)
	Update(ctx context.Context, id int64, req *UpdateSubscriptionRequest) (*Subscription, error)
	ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest, changedBy int64) (*Subscription, error)
	GetStatusHistory(ctx context.Context, id int64) ([]*SubscriptionStatusHistory, error)
	GetActiveByStrategyID(ctx context.Context, strategyID int64) ([]*Subscription, error)
	GetByOfferID(ctx context.Context, offerID int64) ([]*Subscription, error)
	ArchiveByStrategyID(ctx context.Context, strategyID int64, reason string) error
}

type repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) Create(ctx context.Context, req *CreateSubscriptionRequest) (*Subscription, error) {
	query := `
		INSERT INTO subscriptions (investor_user_id, investor_account_id, offer_id, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, investor_user_id, investor_account_id, offer_id, status, created_at, updated_at
	`

	var subscription Subscription
	err := r.db.QueryRowxContext(ctx, query,
		req.InvestorUserID,
		req.InvestorAccountID,
		req.OfferID,
		common.SubscriptionStatusPreparing,
	).StructScan(&subscription)
	if err != nil {
		r.logger.Error("Failed to create subscription",
			zap.Int64("investor_account_id", req.InvestorAccountID),
			zap.Int64("offer_id", req.OfferID),
			zap.Error(err))
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	r.logger.Info("Subscription created",
		zap.Int64("id", subscription.ID),
		zap.Int64("offer_id", subscription.OfferID))

	return &subscription, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*Subscription, error) {
	query := `
		SELECT id, investor_user_id, investor_account_id, offer_id, status, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var subscription Subscription
	err := r.db.GetContext(ctx, &subscription, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get subscription by ID",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("get subscription by id: %w", err)
	}

	return &subscription, nil
}

func (r *repository) List(ctx context.Context, filter *SubscriptionFilter) (*common.PaginatedResult[Subscription], error) {
	filter.SetDefaults()

	var (
		conditions []string
		args       []interface{}
		argIndex   = 1
	)

	if filter.UserID != 0 {
		conditions = append(conditions, fmt.Sprintf("investor_user_id = $%d", argIndex))
		args = append(args, filter.UserID)
		argIndex++
	}

	if filter.OfferID != 0 {
		conditions = append(conditions, fmt.Sprintf("offer_id = $%d", argIndex))
		args = append(args, filter.OfferID)
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

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM subscriptions %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		r.logger.Error("Failed to count subscriptions", zap.Error(err))
		return nil, fmt.Errorf("count subscriptions: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, investor_user_id, investor_account_id, offer_id, status, created_at, updated_at
		FROM subscriptions
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	var subscriptions []Subscription
	err = r.db.SelectContext(ctx, &subscriptions, query, args...)
	if err != nil {
		r.logger.Error("Failed to list subscriptions", zap.Error(err))
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}

	return &common.PaginatedResult[Subscription]{
		Data:       subscriptions,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(filter.Limit))),
	}, nil
}

func (r *repository) Update(ctx context.Context, id int64, req *UpdateSubscriptionRequest) (*Subscription, error) {

	return r.GetByID(ctx, id)
}

func (r *repository) ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest, changedBy int64) (*Subscription, error) {

	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, fmt.Errorf("subscription not found: %d", id)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	updateQuery := `
		UPDATE subscriptions
		SET status = $1, updated_at = now()
		WHERE id = $2
		RETURNING id, investor_user_id, investor_account_id, offer_id, status, created_at, updated_at
	`

	var subscription Subscription
	err = tx.QueryRowxContext(ctx, updateQuery, req.Status, id).StructScan(&subscription)
	if err != nil {
		r.logger.Error("Failed to change subscription status",
			zap.Int64("id", id),
			zap.String("status", string(req.Status)),
			zap.Error(err))
		return nil, fmt.Errorf("change subscription status: %w", err)
	}

	historyQuery := `
		INSERT INTO subscription_status_history (subscription_id, old_status, new_status, reason, changed_by)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.ExecContext(ctx, historyQuery, id, current.Status, req.Status, req.StatusReason, changedBy)
	if err != nil {

		r.logger.Warn("Failed to record status history (table might not exist)",
			zap.Int64("id", id),
			zap.Error(err))
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	r.logger.Info("Subscription status changed",
		zap.Int64("id", subscription.ID),
		zap.String("old_status", string(current.Status)),
		zap.String("new_status", string(subscription.Status)))

	return &subscription, nil
}

func (r *repository) GetStatusHistory(ctx context.Context, id int64) ([]*SubscriptionStatusHistory, error) {
	query := `
		SELECT id, subscription_id, old_status, new_status, reason, changed_by, created_at
		FROM subscription_status_history
		WHERE subscription_id = $1
		ORDER BY created_at DESC
	`

	var history []*SubscriptionStatusHistory
	err := r.db.SelectContext(ctx, &history, query, id)
	if err != nil {
		r.logger.Error("Failed to get subscription status history",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("get subscription status history: %w", err)
	}

	return history, nil
}

func (r *repository) GetActiveByStrategyID(ctx context.Context, strategyID int64) ([]*Subscription, error) {
	query := `
		SELECT s.id, s.investor_user_id, s.investor_account_id, s.offer_id, s.status, s.created_at, s.updated_at
		FROM subscriptions s
		JOIN offers o ON s.offer_id = o.id
		WHERE o.strategy_id = $1 AND s.status = 'active'
		ORDER BY s.created_at DESC
	`

	var subscriptions []*Subscription
	err := r.db.SelectContext(ctx, &subscriptions, query, strategyID)
	if err != nil {
		r.logger.Error("Failed to get active subscriptions by strategy ID",
			zap.Int64("strategy_id", strategyID),
			zap.Error(err))
		return nil, fmt.Errorf("get active subscriptions by strategy id: %w", err)
	}

	return subscriptions, nil
}

func (r *repository) GetByOfferID(ctx context.Context, offerID int64) ([]*Subscription, error) {
	query := `
		SELECT id, investor_user_id, investor_account_id, offer_id, status, created_at, updated_at
		FROM subscriptions
		WHERE offer_id = $1
		ORDER BY created_at DESC
	`

	var subscriptions []*Subscription
	err := r.db.SelectContext(ctx, &subscriptions, query, offerID)
	if err != nil {
		r.logger.Error("Failed to get subscriptions by offer ID",
			zap.Int64("offer_id", offerID),
			zap.Error(err))
		return nil, fmt.Errorf("get subscriptions by offer id: %w", err)
	}

	return subscriptions, nil
}

func (r *repository) ArchiveByStrategyID(ctx context.Context, strategyID int64, reason string) error {
	query := `
		UPDATE subscriptions s
		SET status = 'archived', updated_at = now()
		FROM offers o
		WHERE s.offer_id = o.id
		AND o.strategy_id = $1
		AND s.status = 'active'
	`

	result, err := r.db.ExecContext(ctx, query, strategyID)
	if err != nil {
		r.logger.Error("Failed to archive subscriptions by strategy ID",
			zap.Int64("strategy_id", strategyID),
			zap.String("reason", reason),
			zap.Error(err))
		return fmt.Errorf("archive subscriptions by strategy id: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	r.logger.Info("Subscriptions archived by strategy ID",
		zap.Int64("strategy_id", strategyID),
		zap.Int64("count", rowsAffected),
		zap.String("reason", reason))

	return nil
}
