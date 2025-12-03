package statistics

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type UseCase interface {
	GetStrategyLeaderboard(ctx context.Context, req *LeaderboardRequest) ([]*StrategyLeaderboard, error)
	GetInvestorPortfolio(ctx context.Context, req *InvestorPortfolioRequest) (*InvestorPortfolio, error)
	GetMasterIncome(ctx context.Context, req *MasterIncomeRequest) (*MasterIncome, error)
}

type useCase struct {
	repo   Repository
	logger *zap.Logger
}

func NewUseCase(repo Repository, logger *zap.Logger) UseCase {
	return &useCase{repo: repo, logger: logger}
}

func (u *useCase) GetStrategyLeaderboard(ctx context.Context, req *LeaderboardRequest) ([]*StrategyLeaderboard, error) {

	if req.Limit <= 0 {
		req.Limit = 10
	}

	u.logger.Info("UseCase: Getting strategy leaderboard", zap.Int("limit", req.Limit))

	leaderboard, err := u.repo.GetStrategyLeaderboard(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get strategy leaderboard: %w", err)
	}

	return leaderboard, nil
}

func (u *useCase) GetInvestorPortfolio(ctx context.Context, req *InvestorPortfolioRequest) (*InvestorPortfolio, error) {
	u.logger.Info("UseCase: Getting investor portfolio", zap.Int64("user_id", req.UserID))

	portfolio, err := u.repo.GetInvestorPortfolio(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get investor portfolio: %w", err)
	}

	return portfolio, nil
}

func (u *useCase) GetMasterIncome(ctx context.Context, req *MasterIncomeRequest) (*MasterIncome, error) {
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
