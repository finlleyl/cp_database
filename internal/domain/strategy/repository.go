package strategy

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
	Create(ctx context.Context, req *CreateStrategyRequest) (*Strategy, error)
	GetByID(ctx context.Context, id int64) (*GetStrategyByIDResponse, error)
	GetBaseByID(ctx context.Context, id int64) (*Strategy, error)
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

func NewRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) Create(ctx context.Context, req *CreateStrategyRequest) (*Strategy, error) {
	query := `
		INSERT INTO strategies (master_user_id, master_account_id, title, description, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, master_user_id, master_account_id, title, description, status, created_at, updated_at
	`

	var strategy Strategy
	err := r.db.QueryRowxContext(ctx, query,
		req.UserID,
		req.AccountID,
		req.Nickname,
		req.Summary,
		common.StrategyStatusPreparing,
	).StructScan(&strategy)
	if err != nil {
		r.logger.Error("Failed to create strategy",
			zap.String("nickname", req.Nickname),
			zap.Int64("account_id", req.AccountID),
			zap.Error(err))
		return nil, fmt.Errorf("create strategy: %w", err)
	}

	r.logger.Info("Strategy created",
		zap.Int64("id", strategy.ID),
		zap.String("title", strategy.Title))

	return &strategy, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*GetStrategyByIDResponse, error) {
	query := `SELECT * FROM vw_strategy_performance WHERE id = $1`

	var response GetStrategyByIDResponse
	err := r.db.GetContext(ctx, &response, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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
	var (
		setClauses []string
		args       []interface{}
		argIndex   = 1
	)

	if req.Nickname != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *req.Nickname)
		argIndex++
	}

	if req.Summary != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Summary)
		argIndex++
	}

	if len(setClauses) == 0 {
		strategy, err := r.GetBaseByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return strategy, nil
	}

	setClauses = append(setClauses, "updated_at = now()")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE strategies
		SET %s
		WHERE id = $%d
		RETURNING id, master_user_id, master_account_id, title, description, status, created_at, updated_at
	`, strings.Join(setClauses, ", "), argIndex)

	var strategy Strategy
	err := r.db.QueryRowxContext(ctx, query, args...).StructScan(&strategy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("strategy not found: %d", id)
		}
		r.logger.Error("Failed to update strategy",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("update strategy: %w", err)
	}

	r.logger.Info("Strategy updated", zap.Int64("id", strategy.ID))

	return &strategy, nil
}

func (r *repository) ChangeStatus(ctx context.Context, id int64, req *ChangeStatusRequest) (*Strategy, error) {
	query := `
		UPDATE strategies
		SET status = $1, updated_at = now()
		WHERE id = $2
		RETURNING id, master_user_id, master_account_id, title, description, status, created_at, updated_at
	`

	var strategy Strategy
	err := r.db.QueryRowxContext(ctx, query, req.Status, id).StructScan(&strategy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("strategy not found: %d", id)
		}
		r.logger.Error("Failed to change strategy status",
			zap.Int64("id", id),
			zap.String("status", string(req.Status)),
			zap.Error(err))
		return nil, fmt.Errorf("change strategy status: %w", err)
	}

	r.logger.Info("Strategy status changed",
		zap.Int64("id", strategy.ID),
		zap.String("status", string(strategy.Status)))

	return &strategy, nil
}

func (r *repository) GetByAccountID(ctx context.Context, accountID int64) (*Strategy, error) {
	query := `
		SELECT id, master_user_id, master_account_id, title, description, status, created_at, updated_at
		FROM strategies
		WHERE master_account_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var strategy Strategy
	err := r.db.GetContext(ctx, &strategy, query, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get strategy by account ID",
			zap.Int64("account_id", accountID),
			zap.Error(err))
		return nil, fmt.Errorf("get strategy by account id: %w", err)
	}

	return &strategy, nil
}

func (r *repository) GetActiveByID(ctx context.Context, id int64) (*Strategy, error) {
	query := `
		SELECT id, master_user_id, master_account_id, title, description, status, created_at, updated_at
		FROM strategies
		WHERE id = $1 AND status = 'active'
	`

	var strategy Strategy
	err := r.db.GetContext(ctx, &strategy, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get active strategy by ID",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("get active strategy by id: %w", err)
	}

	return &strategy, nil
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

func (r *repository) GetBaseByID(ctx context.Context, id int64) (*Strategy, error) {
	query := `
		SELECT id, master_user_id, master_account_id, title, description, status, created_at, updated_at
		FROM strategies
		WHERE id = $1
	`

	var strategy Strategy
	err := r.db.GetContext(ctx, &strategy, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get strategy base by ID",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("get strategy base by id: %w", err)
	}

	return &strategy, nil
}
