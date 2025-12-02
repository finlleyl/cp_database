package trade

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Repository defines the interface for trade data operations
type Repository interface {
	Create(ctx context.Context, req *CreateTradeRequest) (*Trade, error)
	GetByID(ctx context.Context, id int64) (*Trade, error)
	List(ctx context.Context, filter *TradeFilter) (*common.PaginatedResult[Trade], error)
	GetByStrategyID(ctx context.Context, strategyID int64, filter *TradeFilter) ([]*Trade, error)
	UpdateProfit(ctx context.Context, id int64, profit float64) error
	CloseTrade(ctx context.Context, id int64, closePrice float64, closeTime time.Time) error
}

// CopiedTradeRepository defines the interface for copied trade data operations
type CopiedTradeRepository interface {
	Create(ctx context.Context, req *CreateCopiedTradeRequest) (*CopiedTrade, error)
	GetByID(ctx context.Context, id int64) (*CopiedTrade, error)
	List(ctx context.Context, filter *CopiedTradeFilter) (*common.PaginatedResult[CopiedTrade], error)
	GetBySubscriptionID(ctx context.Context, subscriptionID int64) ([]*CopiedTrade, error)
	GetByTradeID(ctx context.Context, tradeID int64) ([]*CopiedTrade, error)
	UpdateProfit(ctx context.Context, id int64, profit float64) error
	CloseTrade(ctx context.Context, id int64, closeTime time.Time) error
}

type repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewRepository creates a new trade repository
func NewRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) Create(ctx context.Context, req *CreateTradeRequest) (*Trade, error) {
	query := `
		INSERT INTO trades (strategy_id, master_account_id, symbol, volume_lots, direction, open_time, open_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, strategy_id, master_account_id, symbol, volume_lots, direction, open_time, close_time, open_price, close_price, profit, commission, swap, created_at
	`

	var trade Trade
	err := r.db.QueryRowxContext(ctx, query,
		req.StrategyID,
		req.MasterAccountID,
		req.Symbol,
		req.VolumeLots,
		req.Direction,
		req.OpenTime,
		req.OpenPrice,
	).StructScan(&trade)
	if err != nil {
		r.logger.Error("Failed to create trade",
			zap.Int64("strategy_id", req.StrategyID),
			zap.String("symbol", req.Symbol),
			zap.Error(err))
		return nil, fmt.Errorf("create trade: %w", err)
	}

	r.logger.Info("Trade created",
		zap.Int64("id", trade.ID),
		zap.Int64("strategy_id", trade.StrategyID),
		zap.String("symbol", trade.Symbol))

	return &trade, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*Trade, error) {
	query := `
		SELECT id, strategy_id, master_account_id, symbol, volume_lots, direction, open_time, close_time, open_price, close_price, profit, commission, swap, created_at
		FROM trades
		WHERE id = $1
	`

	var trade Trade
	err := r.db.GetContext(ctx, &trade, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get trade by ID",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("get trade by id: %w", err)
	}

	return &trade, nil
}

func (r *repository) List(ctx context.Context, filter *TradeFilter) (*common.PaginatedResult[Trade], error) {
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

	if !filter.From.IsZero() {
		conditions = append(conditions, fmt.Sprintf("open_time >= $%d", argIndex))
		args = append(args, filter.From)
		argIndex++
	}

	if !filter.To.IsZero() {
		conditions = append(conditions, fmt.Sprintf("open_time <= $%d", argIndex))
		args = append(args, filter.To)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM trades %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		r.logger.Error("Failed to count trades", zap.Error(err))
		return nil, fmt.Errorf("count trades: %w", err)
	}

	// Get data
	query := fmt.Sprintf(`
		SELECT id, strategy_id, master_account_id, symbol, volume_lots, direction, open_time, close_time, open_price, close_price, profit, commission, swap, created_at
		FROM trades
		%s
		ORDER BY open_time DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	var trades []Trade
	err = r.db.SelectContext(ctx, &trades, query, args...)
	if err != nil {
		r.logger.Error("Failed to list trades", zap.Error(err))
		return nil, fmt.Errorf("list trades: %w", err)
	}

	return &common.PaginatedResult[Trade]{
		Data:       trades,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(filter.Limit))),
	}, nil
}

func (r *repository) GetByStrategyID(ctx context.Context, strategyID int64, filter *TradeFilter) ([]*Trade, error) {
	var (
		conditions []string
		args       []interface{}
		argIndex   = 1
	)

	conditions = append(conditions, fmt.Sprintf("strategy_id = $%d", argIndex))
	args = append(args, strategyID)
	argIndex++

	if filter != nil {
		if !filter.From.IsZero() {
			conditions = append(conditions, fmt.Sprintf("open_time >= $%d", argIndex))
			args = append(args, filter.From)
			argIndex++
		}

		if !filter.To.IsZero() {
			conditions = append(conditions, fmt.Sprintf("open_time <= $%d", argIndex))
			args = append(args, filter.To)
			argIndex++
		}
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	query := fmt.Sprintf(`
		SELECT id, strategy_id, master_account_id, symbol, volume_lots, direction, open_time, close_time, open_price, close_price, profit, commission, swap, created_at
		FROM trades
		%s
		ORDER BY open_time DESC
	`, whereClause)

	var trades []*Trade
	err := r.db.SelectContext(ctx, &trades, query, args...)
	if err != nil {
		r.logger.Error("Failed to get trades by strategy ID",
			zap.Int64("strategy_id", strategyID),
			zap.Error(err))
		return nil, fmt.Errorf("get trades by strategy id: %w", err)
	}

	return trades, nil
}

func (r *repository) UpdateProfit(ctx context.Context, id int64, profit float64) error {
	query := `UPDATE trades SET profit = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, profit, id)
	if err != nil {
		r.logger.Error("Failed to update trade profit",
			zap.Int64("id", id),
			zap.Float64("profit", profit),
			zap.Error(err))
		return fmt.Errorf("update trade profit: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("trade not found: %d", id)
	}

	r.logger.Info("Trade profit updated",
		zap.Int64("id", id),
		zap.Float64("profit", profit))

	return nil
}

func (r *repository) CloseTrade(ctx context.Context, id int64, closePrice float64, closeTime time.Time) error {
	query := `UPDATE trades SET close_price = $1, close_time = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, closePrice, closeTime, id)
	if err != nil {
		r.logger.Error("Failed to close trade",
			zap.Int64("id", id),
			zap.Float64("close_price", closePrice),
			zap.Error(err))
		return fmt.Errorf("close trade: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("trade not found: %d", id)
	}

	r.logger.Info("Trade closed",
		zap.Int64("id", id),
		zap.Float64("close_price", closePrice))

	return nil
}

// CopiedTradeRepository implementation

type copiedTradeRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewCopiedTradeRepository creates a new copied trade repository
func NewCopiedTradeRepository(db *sqlx.DB, logger *zap.Logger) CopiedTradeRepository {
	return &copiedTradeRepository{db: db, logger: logger}
}

func (r *copiedTradeRepository) Create(ctx context.Context, req *CreateCopiedTradeRequest) (*CopiedTrade, error) {
	query := `
		INSERT INTO copied_trades (trade_id, subscription_id, investor_account_id, volume_lots, open_time)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, trade_id, subscription_id, investor_account_id, volume_lots, profit, commission, swap, open_time, close_time, created_at
	`

	var copiedTrade CopiedTrade
	err := r.db.QueryRowxContext(ctx, query,
		req.TradeID,
		req.SubscriptionID,
		req.InvestorAccountID,
		req.VolumeLots,
		req.OpenTime,
	).StructScan(&copiedTrade)
	if err != nil {
		r.logger.Error("Failed to create copied trade",
			zap.Int64("trade_id", req.TradeID),
			zap.Int64("subscription_id", req.SubscriptionID),
			zap.Error(err))
		return nil, fmt.Errorf("create copied trade: %w", err)
	}

	r.logger.Info("Copied trade created",
		zap.Int64("id", copiedTrade.ID),
		zap.Int64("trade_id", copiedTrade.TradeID))

	return &copiedTrade, nil
}

func (r *copiedTradeRepository) GetByID(ctx context.Context, id int64) (*CopiedTrade, error) {
	query := `
		SELECT id, trade_id, subscription_id, investor_account_id, volume_lots, profit, commission, swap, open_time, close_time, created_at
		FROM copied_trades
		WHERE id = $1
	`

	var copiedTrade CopiedTrade
	err := r.db.GetContext(ctx, &copiedTrade, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get copied trade by ID",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("get copied trade by id: %w", err)
	}

	return &copiedTrade, nil
}

func (r *copiedTradeRepository) List(ctx context.Context, filter *CopiedTradeFilter) (*common.PaginatedResult[CopiedTrade], error) {
	filter.SetDefaults()

	var (
		conditions []string
		args       []interface{}
		argIndex   = 1
	)

	if filter.SubscriptionID != 0 {
		conditions = append(conditions, fmt.Sprintf("subscription_id = $%d", argIndex))
		args = append(args, filter.SubscriptionID)
		argIndex++
	}

	if filter.TradeID != 0 {
		conditions = append(conditions, fmt.Sprintf("trade_id = $%d", argIndex))
		args = append(args, filter.TradeID)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM copied_trades %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		r.logger.Error("Failed to count copied trades", zap.Error(err))
		return nil, fmt.Errorf("count copied trades: %w", err)
	}

	// Get data
	query := fmt.Sprintf(`
		SELECT id, trade_id, subscription_id, investor_account_id, volume_lots, profit, commission, swap, open_time, close_time, created_at
		FROM copied_trades
		%s
		ORDER BY open_time DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	var copiedTrades []CopiedTrade
	err = r.db.SelectContext(ctx, &copiedTrades, query, args...)
	if err != nil {
		r.logger.Error("Failed to list copied trades", zap.Error(err))
		return nil, fmt.Errorf("list copied trades: %w", err)
	}

	return &common.PaginatedResult[CopiedTrade]{
		Data:       copiedTrades,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(filter.Limit))),
	}, nil
}

func (r *copiedTradeRepository) GetBySubscriptionID(ctx context.Context, subscriptionID int64) ([]*CopiedTrade, error) {
	query := `
		SELECT id, trade_id, subscription_id, investor_account_id, volume_lots, profit, commission, swap, open_time, close_time, created_at
		FROM copied_trades
		WHERE subscription_id = $1
		ORDER BY open_time DESC
	`

	var copiedTrades []*CopiedTrade
	err := r.db.SelectContext(ctx, &copiedTrades, query, subscriptionID)
	if err != nil {
		r.logger.Error("Failed to get copied trades by subscription ID",
			zap.Int64("subscription_id", subscriptionID),
			zap.Error(err))
		return nil, fmt.Errorf("get copied trades by subscription id: %w", err)
	}

	return copiedTrades, nil
}

func (r *copiedTradeRepository) GetByTradeID(ctx context.Context, tradeID int64) ([]*CopiedTrade, error) {
	query := `
		SELECT id, trade_id, subscription_id, investor_account_id, volume_lots, profit, commission, swap, open_time, close_time, created_at
		FROM copied_trades
		WHERE trade_id = $1
		ORDER BY created_at DESC
	`

	var copiedTrades []*CopiedTrade
	err := r.db.SelectContext(ctx, &copiedTrades, query, tradeID)
	if err != nil {
		r.logger.Error("Failed to get copied trades by trade ID",
			zap.Int64("trade_id", tradeID),
			zap.Error(err))
		return nil, fmt.Errorf("get copied trades by trade id: %w", err)
	}

	return copiedTrades, nil
}

func (r *copiedTradeRepository) UpdateProfit(ctx context.Context, id int64, profit float64) error {
	query := `UPDATE copied_trades SET profit = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, profit, id)
	if err != nil {
		r.logger.Error("Failed to update copied trade profit",
			zap.Int64("id", id),
			zap.Float64("profit", profit),
			zap.Error(err))
		return fmt.Errorf("update copied trade profit: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("copied trade not found: %d", id)
	}

	r.logger.Info("Copied trade profit updated",
		zap.Int64("id", id),
		zap.Float64("profit", profit))

	return nil
}

func (r *copiedTradeRepository) CloseTrade(ctx context.Context, id int64, closeTime time.Time) error {
	query := `UPDATE copied_trades SET close_time = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, closeTime, id)
	if err != nil {
		r.logger.Error("Failed to close copied trade",
			zap.Int64("id", id),
			zap.Error(err))
		return fmt.Errorf("close copied trade: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("copied trade not found: %d", id)
	}

	r.logger.Info("Copied trade closed", zap.Int64("id", id))

	return nil
}
