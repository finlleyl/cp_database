package statistics

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Repository defines the interface for statistics data operations
type Repository interface {
	GetStrategyLeaderboard(ctx context.Context, req *LeaderboardRequest) ([]*StrategyLeaderboard, error)
	GetInvestorPortfolio(ctx context.Context, req *InvestorPortfolioRequest) (*InvestorPortfolio, error)
	GetMasterIncome(ctx context.Context, req *MasterIncomeRequest) (*MasterIncome, error)
	CreateCommission(ctx context.Context, req *CreateCommissionRequest) (*Commission, error)
	GetCommissionsBySubscriptionID(ctx context.Context, subscriptionID int64) ([]*Commission, error)
}

type repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewRepository creates a new statistics repository
func NewRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) GetStrategyLeaderboard(ctx context.Context, req *LeaderboardRequest) ([]*StrategyLeaderboard, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	r.logger.Info("Getting strategy leaderboard", zap.Int("limit", req.Limit))

	query := `SELECT strategy_id, title, total_profit, total_commissions, active_subscriptions 
			  FROM fn_get_strategy_leaderboard($1)`

	var leaderboard []*StrategyLeaderboard
	if err := r.db.SelectContext(ctx, &leaderboard, query, req.Limit); err != nil {
		r.logger.Error("Failed to get strategy leaderboard", zap.Error(err))
		return nil, fmt.Errorf("get strategy leaderboard: %w", err)
	}

	return leaderboard, nil
}

func (r *repository) GetInvestorPortfolio(ctx context.Context, req *InvestorPortfolioRequest) (*InvestorPortfolio, error) {
	r.logger.Info("Getting investor portfolio", zap.Int64("user_id", req.UserID))

	query := `SELECT subscription_id, strategy_id, strategy_title, total_profit, copied_trades_count 
			  FROM fn_get_investor_portfolio($1)`

	var items []PortfolioItem
	if err := r.db.SelectContext(ctx, &items, query, req.UserID); err != nil {
		r.logger.Error("Failed to get investor portfolio",
			zap.Int64("user_id", req.UserID),
			zap.Error(err))
		return nil, fmt.Errorf("get investor portfolio: %w", err)
	}

	return &InvestorPortfolio{
		UserID:        req.UserID,
		Subscriptions: items,
	}, nil
}

func (r *repository) GetMasterIncome(ctx context.Context, req *MasterIncomeRequest) (*MasterIncome, error) {
	r.logger.Info("Getting master income",
		zap.Int64("user_id", req.UserID),
		zap.Time("from", req.From),
		zap.Time("to", req.To))

	// Query to get total commissions by type for a master user
	// This joins commissions through subscriptions -> offers -> strategies to find master's income
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN c.type = 'performance' THEN c.amount ELSE 0 END), 0) as performance_fees,
			COALESCE(SUM(CASE WHEN c.type = 'management' THEN c.amount ELSE 0 END), 0) as management_fees,
			COALESCE(SUM(CASE WHEN c.type = 'registration' THEN c.amount ELSE 0 END), 0) as registration_fees
		FROM commissions c
		JOIN subscriptions s ON c.subscription_id = s.id
		JOIN offers o ON s.offer_id = o.id
		JOIN strategies st ON o.strategy_id = st.id
		WHERE st.master_user_id = $1
	`

	args := []interface{}{req.UserID}
	argIndex := 2

	if !req.From.IsZero() {
		query += fmt.Sprintf(" AND c.created_at >= $%d", argIndex)
		args = append(args, req.From)
		argIndex++
	}

	if !req.To.IsZero() {
		query += fmt.Sprintf(" AND c.created_at <= $%d", argIndex)
		args = append(args, req.To)
	}

	var result struct {
		PerformanceFees  float64 `db:"performance_fees"`
		ManagementFees   float64 `db:"management_fees"`
		RegistrationFees float64 `db:"registration_fees"`
	}

	err := r.db.GetContext(ctx, &result, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return &MasterIncome{
				UserID:           req.UserID,
				TotalIncome:      0,
				PerformanceFees:  0,
				ManagementFees:   0,
				RegistrationFees: 0,
			}, nil
		}
		r.logger.Error("Failed to get master income",
			zap.Int64("user_id", req.UserID),
			zap.Error(err))
		return nil, fmt.Errorf("get master income: %w", err)
	}

	return &MasterIncome{
		UserID:           req.UserID,
		TotalIncome:      result.PerformanceFees + result.ManagementFees + result.RegistrationFees,
		PerformanceFees:  result.PerformanceFees,
		ManagementFees:   result.ManagementFees,
		RegistrationFees: result.RegistrationFees,
	}, nil
}

func (r *repository) CreateCommission(ctx context.Context, req *CreateCommissionRequest) (*Commission, error) {
	query := `
		INSERT INTO commissions (subscription_id, type, amount, period_from, period_to)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, subscription_id, type, amount, period_from, period_to, created_at
	`

	var commission Commission
	err := r.db.QueryRowxContext(ctx, query,
		req.SubscriptionID,
		req.Type,
		req.Amount,
		req.PeriodFrom,
		req.PeriodTo,
	).StructScan(&commission)
	if err != nil {
		r.logger.Error("Failed to create commission",
			zap.Int64("subscription_id", req.SubscriptionID),
			zap.String("type", string(req.Type)),
			zap.Float64("amount", req.Amount),
			zap.Error(err))
		return nil, fmt.Errorf("create commission: %w", err)
	}

	r.logger.Info("Commission created",
		zap.Int64("id", commission.ID),
		zap.Int64("subscription_id", commission.SubscriptionID),
		zap.String("type", string(commission.Type)))

	return &commission, nil
}

func (r *repository) GetCommissionsBySubscriptionID(ctx context.Context, subscriptionID int64) ([]*Commission, error) {
	query := `
		SELECT id, subscription_id, type, amount, period_from, period_to, created_at
		FROM commissions
		WHERE subscription_id = $1
		ORDER BY created_at DESC
	`

	var commissions []*Commission
	err := r.db.SelectContext(ctx, &commissions, query, subscriptionID)
	if err != nil {
		r.logger.Error("Failed to get commissions by subscription ID",
			zap.Int64("subscription_id", subscriptionID),
			zap.Error(err))
		return nil, fmt.Errorf("get commissions by subscription id: %w", err)
	}

	return commissions, nil
}
