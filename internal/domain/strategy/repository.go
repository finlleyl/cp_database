package strategy

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Repository defines the interface for strategy data operations
type Repository interface {
	Create(ctx context.Context, req *CreateStrategyRequest) (*Strategy, error)
	GetByID(ctx context.Context, id int64) (*GetStrategyByIDResponse, error)
	List(ctx context.Context, filter *StrategyFilter) (*common.PaginatedResult[GetStrategyByIDResponse], error)
	Update(ctx context.Context, id int64, req *UpdateStrategyRequest) (*Strategy, error)
	ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest) (*Strategy, error)
	GetByAccountID(ctx context.Context, accountID int64) (*Strategy, error)
	GetActiveByID(ctx context.Context, id int64) (*Strategy, error)
	GetSummary(ctx context.Context, id int64) (*StrategySummary, error)
}

type repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewRepository creates a new strategy repository
func NewRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) Create(ctx context.Context, req *CreateStrategyRequest) (*Strategy, error) {
	// TODO: Implement strategy creation with initial status = preparing
	r.logger.Info("Creating strategy", zap.String("nickname", req.Nickname), zap.Int64("account_id", req.AccountID))
	return nil, fmt.Errorf("not implemented")
}

func (r *repository) GetByID(ctx context.Context, id int64) (*GetStrategyByIDResponse, error) {
	query := `
		SELECT * from vw_strategy_performance where id = $1
	`
	var response GetStrategyByIDResponse
	err := r.db.GetContext(ctx, &response, query, id)
	if err != nil {
		r.logger.Error("Failed to get strategy by ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get strategy by ID: %w", err)
	}
	return &response, nil
}

func (r *repository) List(ctx context.Context, filter *StrategyFilter) (*common.PaginatedResult[GetStrategyByIDResponse], error) {
	var (
		where  []string
		args   []interface{}
		argPos = 1
	)

	if filter.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", argPos))
		args = append(args, filter.Status)
		argPos++
	}

	if filter.MinROI != nil {
		where = append(where, fmt.Sprintf("roi >= $%d", argPos))
		args = append(args, *filter.MinROI)
		argPos++
	}

	if filter.MaxDrawdownPct != nil {
		where = append(where, fmt.Sprintf("max_drawdown_pct <= $%d", argPos))
		args = append(args, *filter.MaxDrawdownPct)
		argPos++
	}

	if filter.RiskScore != nil {
		where = append(where, fmt.Sprintf("risk_score = $%d", argPos))
		args = append(args, *filter.RiskScore)
		argPos++
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	countQuery := `SELECT COUNT(*) FROM vw_strategy_performance ` + whereSQL

	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, fmt.Errorf("count failed: %w", err)
	}

	filter.Pagination.SetDefaults()

	mainQuery := fmt.Sprintf(`
		SELECT *
		FROM vw_strategy_performance
		%s
		ORDER BY total_profit DESC
		LIMIT %d OFFSET %d
	`, whereSQL, filter.Limit, filter.Offset)

	var items []GetStrategyByIDResponse
	if err := r.db.SelectContext(ctx, &items, mainQuery, args...); err != nil {
		return nil, fmt.Errorf("list failed: %w", err)
	}

	return &common.PaginatedResult[GetStrategyByIDResponse]{
		Data:       items,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(filter.Limit))),
	}, nil
}

func (r *repository) Update(ctx context.Context, id int64, req *UpdateStrategyRequest) (*Strategy, error) {
	// TODO: Implement strategy update
	r.logger.Info("Updating strategy", zap.Int64("id", id))
	return nil, fmt.Errorf("not implemented")
}

func (r *repository) ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest) (*Strategy, error) {
	// TODO: Implement strategy status change with status_reason
	r.logger.Info("Changing strategy status", zap.Int64("id", id), zap.String("status", string(req.Status)))
	return nil, fmt.Errorf("not implemented")
}

func (r *repository) GetByAccountID(ctx context.Context, accountID int64) (*Strategy, error) {
	// TODO: Implement get strategy by account ID
	r.logger.Info("Getting strategy by account ID", zap.Int64("account_id", accountID))
	return nil, fmt.Errorf("not implemented")
}

func (r *repository) GetActiveByID(ctx context.Context, id int64) (*Strategy, error) {
	// TODO: Implement get active strategy by UUID
	r.logger.Info("Getting active strategy by ID", zap.Int64("id", id))
	return nil, fmt.Errorf("not implemented")
}

func (r *repository) GetSummary(ctx context.Context, id int64) (*StrategySummary, error) {
	r.logger.Info("Getting strategy summary", zap.Int64("id", id))

	query := `SELECT fn_get_strategy_total_profit($1) as total_profit`

	var totalProfit float64
	if err := r.db.GetContext(ctx, &totalProfit, query, id); err != nil {
		return nil, fmt.Errorf("get strategy total profit: %w", err)
	}

	return &StrategySummary{
		StrategyID:  id,
		TotalProfit: totalProfit,
	}, nil
}
