package trade

import (
	"context"
	"fmt"

	"github.com/finlleyl/cp_database/internal/domain/audit"
	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/finlleyl/cp_database/internal/domain/statistics"
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
	statisticsRepo   statistics.Repository
	auditRepo        audit.Repository
	logger           *zap.Logger
}

// NewUseCase creates a new trade use case
func NewUseCase(
	repo Repository,
	copiedTradeRepo CopiedTradeRepository,
	subscriptionRepo subscription.Repository,
	statisticsRepo statistics.Repository,
	auditRepo audit.Repository,
	logger *zap.Logger,
) UseCase {
	return &useCase{
		repo:             repo,
		copiedTradeRepo:  copiedTradeRepo,
		subscriptionRepo: subscriptionRepo,
		statisticsRepo:   statisticsRepo,
		auditRepo:        auditRepo,
		logger:           logger,
	}
}

func (u *useCase) Create(ctx context.Context, req *CreateTradeRequest) (*Trade, error) {
	// TODO: Implement trade creation business logic
	// Business flow:
	// 1. Validate strategy exists and is active
	// 2. Validate account belongs to strategy
	// 3. Create trade
	// 4. Trigger copy to all active subscriptions (async or sync based on config)
	// 5. Update account statistics
	// 6. Create audit log
	u.logger.Info("UseCase: Creating trade",
		zap.String("strategy_uuid", req.StrategyUUID.String()),
		zap.String("symbol", req.Symbol),
		zap.String("type", string(req.Type)),
		zap.Float64("volume", req.Volume))

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
	// TODO: Implement get trade by ID business logic
	u.logger.Info("UseCase: Getting trade by ID", zap.Int64("id", id))
	return u.repo.GetByID(ctx, id)
}

func (u *useCase) List(ctx context.Context, filter *TradeFilter) (*common.PaginatedResult[Trade], error) {
	// TODO: Implement trade listing business logic
	// Supports filtering by strategy_uuid and time range
	filter.SetDefaults()
	u.logger.Info("UseCase: Listing trades", zap.Any("filter", filter))
	return u.repo.List(ctx, filter)
}

func (u *useCase) CopyTrade(ctx context.Context, tradeID int64, req *CopyTradeRequest) ([]*CopiedTrade, error) {
	// TODO: Implement trade copying business logic
	// Business flow:
	// 1. Get original trade
	// 2. Get all active subscriptions for the strategy (or specific ones from request)
	// 3. For each subscription:
	//    a. Apply filter rules (allowed/blocked symbols, lot size limits)
	//    b. Calculate copy ratio
	//    c. Create copied trade
	//    d. Update investor account statistics
	//    e. Calculate performance fee and create commission record
	// 4. Return list of created copied trades
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
		// TODO: Apply filter rules from subscription
		// TODO: Calculate volume based on copy ratio
		copyRatio := 1.0 // Default, should be read from subscription config

		copiedTrade := &CopiedTrade{
			OriginalTradeID: trade.ID,
			SubscriptionID:  sub.ID,
			InvestorAccount: sub.InvestorAccountID,
			Symbol:          trade.Symbol,
			Type:            trade.Type,
			Volume:          trade.Volume * copyRatio,
			CopyRatio:       copyRatio,
			OpenPrice:       trade.OpenPrice,
			StopLoss:        trade.StopLoss,
			TakeProfit:      trade.TakeProfit,
			Status:          TradeStatusOpen,
			OpenedAt:        trade.OpenedAt,
		}

		created, err := u.copiedTradeRepo.Create(ctx, copiedTrade)
		if err != nil {
			u.logger.Error("Failed to create copied trade",
				zap.Error(err),
				zap.Int64("subscription_id", sub.ID))
			continue
		}

		copiedTrades = append(copiedTrades, created)

		// TODO: Update account statistics
		// TODO: Calculate and record performance fee
	}

	u.logger.Info("Trade copied",
		zap.Int64("trade_id", tradeID),
		zap.Int("copied_count", len(copiedTrades)))

	return copiedTrades, nil
}

func (u *useCase) ListCopiedTrades(ctx context.Context, filter *CopiedTradeFilter) (*common.PaginatedResult[CopiedTrade], error) {
	// TODO: Implement copied trade listing business logic
	filter.SetDefaults()
	u.logger.Info("UseCase: Listing copied trades", zap.Any("filter", filter))
	return u.copiedTradeRepo.List(ctx, filter)
}
