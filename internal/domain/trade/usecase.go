package trade

import (
	"context"
	"fmt"
	"time"

	"github.com/finlleyl/cp_database/internal/domain/audit"
	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/finlleyl/cp_database/internal/domain/subscription"
	"go.uber.org/zap"
)

// UseCase defines the interface for trade business logic
type UseCase interface {
	Create(ctx context.Context, req *CreateTradeRequest) (*Trade, error)
	GetByID(ctx context.Context, id int64) (*Trade, error)
	List(ctx context.Context, filter *TradeFilter) (*common.PaginatedResult[Trade], error)
	CopyTrade(ctx context.Context, tradeID int64, req *CopyTradeRequest) ([]*CopiedTrade, error)
	ListCopiedTrades(ctx context.Context, filter *CopiedTradeFilter) (*common.PaginatedResult[CopiedTrade], error)
}

type useCase struct {
	repo             Repository
	copiedTradeRepo  CopiedTradeRepository
	subscriptionRepo subscription.Repository
	auditRepo        audit.Repository
	logger           *zap.Logger
}

// NewUseCase creates a new trade use case
func NewUseCase(
	repo Repository,
	copiedTradeRepo CopiedTradeRepository,
	subscriptionRepo subscription.Repository,
	auditRepo audit.Repository,
	logger *zap.Logger,
) UseCase {
	return &useCase{
		repo:             repo,
		copiedTradeRepo:  copiedTradeRepo,
		subscriptionRepo: subscriptionRepo,
		auditRepo:        auditRepo,
		logger:           logger,
	}
}

func (u *useCase) Create(ctx context.Context, req *CreateTradeRequest) (*Trade, error) {
	u.logger.Info("UseCase: Creating trade",
		zap.Int64("strategy_id", req.StrategyID),
		zap.String("symbol", req.Symbol),
		zap.String("direction", string(req.Direction)),
		zap.Float64("volume_lots", req.VolumeLots))

	trade, err := u.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create trade: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeTrade,
		EntityID:   trade.ID,
		Action:     audit.AuditActionCreate,
		NewValue:   trade,
	})

	return trade, nil
}

func (u *useCase) GetByID(ctx context.Context, id int64) (*Trade, error) {
	u.logger.Info("UseCase: Getting trade by ID", zap.Int64("id", id))
	return u.repo.GetByID(ctx, id)
}

func (u *useCase) List(ctx context.Context, filter *TradeFilter) (*common.PaginatedResult[Trade], error) {
	filter.SetDefaults()
	u.logger.Info("UseCase: Listing trades", zap.Any("filter", filter))
	return u.repo.List(ctx, filter)
}

func (u *useCase) CopyTrade(ctx context.Context, tradeID int64, req *CopyTradeRequest) ([]*CopiedTrade, error) {
	u.logger.Info("UseCase: Copying trade",
		zap.Int64("trade_id", tradeID),
		zap.Any("subscription_ids", req.SubscriptionIDs))

	trade, err := u.repo.GetByID(ctx, tradeID)
	if err != nil {
		return nil, fmt.Errorf("get trade: %w", err)
	}
	if trade == nil {
		return nil, fmt.Errorf("trade not found")
	}

	var subscriptions []*subscription.Subscription
	if len(req.SubscriptionIDs) > 0 {
		// Get specific subscriptions
		for _, subID := range req.SubscriptionIDs {
			sub, err := u.subscriptionRepo.GetByID(ctx, subID)
			if err != nil {
				u.logger.Warn("Failed to get subscription", zap.Error(err))
				continue
			}
			if sub != nil && sub.Status == common.SubscriptionStatusActive {
				subscriptions = append(subscriptions, sub)
			}
		}
	} else {
		// Get all active subscriptions for the strategy
		subscriptions, err = u.subscriptionRepo.GetActiveByStrategyID(ctx, trade.StrategyID)
		if err != nil {
			return nil, fmt.Errorf("get active subscriptions: %w", err)
		}
	}

	var copiedTrades []*CopiedTrade
	for _, sub := range subscriptions {
		copyReq := &CreateCopiedTradeRequest{
			TradeID:           trade.ID,
			SubscriptionID:    sub.ID,
			InvestorAccountID: sub.InvestorAccountID,
			VolumeLots:        trade.VolumeLots,
			OpenTime:          time.Now(),
		}

		created, err := u.copiedTradeRepo.Create(ctx, copyReq)
		if err != nil {
			u.logger.Error("Failed to create copied trade",
				zap.Error(err),
				zap.Int64("subscription_id", sub.ID))
			continue
		}

		copiedTrades = append(copiedTrades, created)
	}

	u.logger.Info("Trade copied",
		zap.Int64("trade_id", tradeID),
		zap.Int("copied_count", len(copiedTrades)))

	return copiedTrades, nil
}

func (u *useCase) ListCopiedTrades(ctx context.Context, filter *CopiedTradeFilter) (*common.PaginatedResult[CopiedTrade], error) {
	filter.SetDefaults()
	u.logger.Info("UseCase: Listing copied trades", zap.Any("filter", filter))
	return u.copiedTradeRepo.List(ctx, filter)
}
