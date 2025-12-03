package account

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
	Create(ctx context.Context, req *CreateAccountRequest) (*Account, error)
	GetByID(ctx context.Context, id int64) (*Account, error)
	List(ctx context.Context, filter *AccountFilter) (*common.PaginatedResult[Account], error)
	Update(ctx context.Context, id int64, req *UpdateAccountRequest) (*Account, error)
	Delete(ctx context.Context, id int64) error
	GetByUserID(ctx context.Context, userID int64) ([]*Account, error)
}

type repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) Create(ctx context.Context, req *CreateAccountRequest) (*Account, error) {
	query := `
		INSERT INTO accounts (user_id, name, account_type, currency)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, name, account_type, currency, created_at, updated_at
	`

	var account Account
	err := r.db.QueryRowxContext(ctx, query,
		req.UserID,
		req.Name,
		req.AccountType,
		req.Currency,
	).StructScan(&account)
	if err != nil {
		r.logger.Error("Failed to create account",
			zap.Int64("user_id", req.UserID),
			zap.Error(err))
		return nil, fmt.Errorf("create account: %w", err)
	}

	r.logger.Info("Account created",
		zap.Int64("id", account.ID),
		zap.Int64("user_id", account.UserID))

	return &account, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*Account, error) {
	query := `
		SELECT id, user_id, name, account_type, currency, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`

	var account Account
	err := r.db.GetContext(ctx, &account, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get account by ID",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("get account by id: %w", err)
	}

	return &account, nil
}

func (r *repository) List(ctx context.Context, filter *AccountFilter) (*common.PaginatedResult[Account], error) {
	filter.SetDefaults()

	var (
		conditions []string
		args       []interface{}
		argIndex   = 1
	)

	if filter.UserID != 0 {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, filter.UserID)
		argIndex++
	}

	if filter.AccountType != "" {
		conditions = append(conditions, fmt.Sprintf("account_type = $%d", argIndex))
		args = append(args, filter.AccountType)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM accounts %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		r.logger.Error("Failed to count accounts", zap.Error(err))
		return nil, fmt.Errorf("count accounts: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, name, account_type, currency, created_at, updated_at
		FROM accounts
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	var accounts []Account
	err = r.db.SelectContext(ctx, &accounts, query, args...)
	if err != nil {
		r.logger.Error("Failed to list accounts", zap.Error(err))
		return nil, fmt.Errorf("list accounts: %w", err)
	}

	return &common.PaginatedResult[Account]{
		Data:       accounts,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(filter.Limit))),
	}, nil
}

func (r *repository) Update(ctx context.Context, id int64, req *UpdateAccountRequest) (*Account, error) {
	var (
		setClauses []string
		args       []interface{}
		argIndex   = 1
	)

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}

	if req.AccountType != nil {
		setClauses = append(setClauses, fmt.Sprintf("account_type = $%d", argIndex))
		args = append(args, *req.AccountType)
		argIndex++
	}

	if req.Currency != nil {
		setClauses = append(setClauses, fmt.Sprintf("currency = $%d", argIndex))
		args = append(args, *req.Currency)
		argIndex++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	setClauses = append(setClauses, "updated_at = now()")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE accounts
		SET %s
		WHERE id = $%d
		RETURNING id, user_id, name, account_type, currency, created_at, updated_at
	`, strings.Join(setClauses, ", "), argIndex)

	var account Account
	err := r.db.QueryRowxContext(ctx, query, args...).StructScan(&account)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found: %d", id)
		}
		r.logger.Error("Failed to update account",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("update account: %w", err)
	}

	r.logger.Info("Account updated", zap.Int64("id", account.ID))

	return &account, nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {

	var strategyCount int
	err := r.db.GetContext(ctx, &strategyCount,
		"SELECT COUNT(*) FROM strategies WHERE master_account_id = $1", id)
	if err != nil {
		r.logger.Error("Failed to check account dependencies",
			zap.Int64("id", id),
			zap.Error(err))
		return fmt.Errorf("check dependencies: %w", err)
	}

	if strategyCount > 0 {
		return fmt.Errorf("cannot delete account: has %d active strategies", strategyCount)
	}

	query := `DELETE FROM accounts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete account",
			zap.Int64("id", id),
			zap.Error(err))
		return fmt.Errorf("delete account: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("account not found: %d", id)
	}

	r.logger.Info("Account deleted", zap.Int64("id", id))

	return nil
}

func (r *repository) GetByUserID(ctx context.Context, userID int64) ([]*Account, error) {
	query := `
		SELECT id, user_id, name, account_type, currency, created_at, updated_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var accounts []*Account
	err := r.db.SelectContext(ctx, &accounts, query, userID)
	if err != nil {
		r.logger.Error("Failed to get accounts by user ID",
			zap.Int64("user_id", userID),
			zap.Error(err))
		return nil, fmt.Errorf("get accounts by user id: %w", err)
	}

	return accounts, nil
}
