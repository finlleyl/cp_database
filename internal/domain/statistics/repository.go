package statistics

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Repository defines the interface for statistics data operations
type Repository interface {
	GetStrategyLeaderboard(ctx context.Context, req *LeaderboardRequest) ([]*StrategyLeaderboard, error)
	GetInvestorPortfolio(ctx context.Context, req *InvestorPortfolioRequest) (*InvestorPortfolio, error)
	GetMasterIncome(ctx context.Context, req *MasterIncomeRequest) (*MasterIncome, error)
	GetAccountStatistics(ctx context.Context, req *AccountStatisticsRequest) (*AccountStatistics, error)
	UpdateAccountStatistics(ctx context.Context, accountID int64) error
	CreateCommission(ctx context.Context, commission *Commission) (*Commission, error)
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
	r.logger.Info("Getting strategy leaderboard", zap.Int("limit", req.Limit))

	query := `SELECT strategy_id, title, total_profit, total_commissions, active_subscriptions 
			  FROM fn_get_strategy_leaderboard($1)`

	var leaderboard []*StrategyLeaderboard
	if err := r.db.SelectContext(ctx, &leaderboard, query, req.Limit); err != nil {
		return nil, fmt.Errorf("select strategy leaderboard: %w", err)
	}

	return leaderboard, nil
}

func (r *repository) GetInvestorPortfolio(ctx context.Context, req *InvestorPortfolioRequest) (*InvestorPortfolio, error) {
	r.logger.Info("Getting investor portfolio", zap.Int64("user_id", req.UserID))

	query := `SELECT subscription_id, strategy_id, strategy_title, total_profit, copied_trades_count 
			  FROM fn_get_investor_portfolio($1)`

	var items []PortfolioItem
	if err := r.db.SelectContext(ctx, &items, query, req.UserID); err != nil {
		return nil, fmt.Errorf("select investor portfolio: %w", err)
	}

	return &InvestorPortfolio{
		UserID:        req.UserID,
		Subscriptions: items,
	}, nil
}

func (r *repository) GetMasterIncome(ctx context.Context, req *MasterIncomeRequest) (*MasterIncome, error) {
	// TODO: Implement using fn_get_master_income database function
	r.logger.Info("Getting master income",
		zap.Int64("user_id", req.UserID),
		zap.Time("from", req.From),
		zap.Time("to", req.To))
	return nil, fmt.Errorf("not implemented")
}

func (r *repository) GetAccountStatistics(ctx context.Context, req *AccountStatisticsRequest) (*AccountStatistics, error) {
	// TODO: Implement get account statistics from account_statistics table
	r.logger.Info("Getting account statistics",
		zap.Int64("account_id", req.AccountID),
		zap.String("period", string(req.Period)))
	return nil, fmt.Errorf("not implemented")
}

func (r *repository) UpdateAccountStatistics(ctx context.Context, accountID int64) error {
	// TODO: Implement account statistics recalculation
	r.logger.Info("Updating account statistics", zap.Int64("account_id", accountID))
	return fmt.Errorf("not implemented")
}

func (r *repository) CreateCommission(ctx context.Context, commission *Commission) (*Commission, error) {
	// TODO: Implement commission creation
	r.logger.Info("Creating commission",
		zap.String("subscription_uuid", commission.SubscriptionUUID.String()),
		zap.String("type", string(commission.Type)),
		zap.Float64("amount", commission.Amount))
	return nil, fmt.Errorf("not implemented")
}
