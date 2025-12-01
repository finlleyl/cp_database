package statistics

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// UseCase defines the interface for statistics business logic
type UseCase interface {
	GetStrategyLeaderboard(ctx context.Context, req *LeaderboardRequest) ([]*StrategyLeaderboard, error)
	GetInvestorPortfolio(ctx context.Context, req *InvestorPortfolioRequest) (*InvestorPortfolio, error)
	GetMasterIncome(ctx context.Context, req *MasterIncomeRequest) (*MasterIncome, error)
	GetAccountStatistics(ctx context.Context, req *AccountStatisticsRequest) (*AccountStatistics, error)
}

type useCase struct {
	repo   Repository
	logger *zap.Logger
}

// NewUseCase creates a new statistics use case
func NewUseCase(repo Repository, logger *zap.Logger) UseCase {
	return &useCase{repo: repo, logger: logger}
}

func (u *useCase) GetStrategyLeaderboard(ctx context.Context, req *LeaderboardRequest) ([]*StrategyLeaderboard, error) {
	// TODO: Implement strategy leaderboard business logic
	// Uses fn_get_strategy_leaderboard database function
	// Default period: month, default limit: 20
	if req.Period == "" {
		req.Period = PeriodMonth
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	u.logger.Info("UseCase: Getting strategy leaderboard",
		zap.String("period", string(req.Period)),
		zap.Int("limit", req.Limit))

	leaderboard, err := u.repo.GetStrategyLeaderboard(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get strategy leaderboard: %w", err)
	}

	return leaderboard, nil
}

func (u *useCase) GetInvestorPortfolio(ctx context.Context, req *InvestorPortfolioRequest) (*InvestorPortfolio, error) {
	// TODO: Implement investor portfolio business logic
	// Uses fn_get_investor_portfolio database function
	// Returns aggregated portfolio data with subscriptions breakdown
	u.logger.Info("UseCase: Getting investor portfolio", zap.Int64("user_id", req.UserID))

	portfolio, err := u.repo.GetInvestorPortfolio(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get investor portfolio: %w", err)
	}

	return portfolio, nil
}

func (u *useCase) GetMasterIncome(ctx context.Context, req *MasterIncomeRequest) (*MasterIncome, error) {
	// TODO: Implement master income business logic
	// Uses fn_get_master_income database function
	// Returns income breakdown by fee type and strategy
	u.logger.Info("UseCase: Getting master income",
		zap.Int64("user_id", req.UserID),
		zap.Time("from", req.From),
		zap.Time("to", req.To))

	income, err := u.repo.GetMasterIncome(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get master income: %w", err)
	}

	return income, nil
}

func (u *useCase) GetAccountStatistics(ctx context.Context, req *AccountStatisticsRequest) (*AccountStatistics, error) {
	// TODO: Implement account statistics business logic
	// Reads from account_statistics table
	// Default period: all
	if req.Period == "" {
		req.Period = PeriodAll
	}

	u.logger.Info("UseCase: Getting account statistics",
		zap.Int64("account_id", req.AccountID),
		zap.String("period", string(req.Period)))

	stats, err := u.repo.GetAccountStatistics(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get account statistics: %w", err)
	}

	return stats, nil
}
