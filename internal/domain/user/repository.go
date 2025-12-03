package user

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
	Create(ctx context.Context, req *CreateUserRequest) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	List(ctx context.Context, filter *UserFilter) (*common.PaginatedResult[User], error)
	Update(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error)
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &repository{db: db, logger: logger}
}

func (r *repository) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
	query := `
		INSERT INTO users (name, email, role)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, role, created_at, updated_at
	`

	var user User
	err := r.db.QueryRowxContext(ctx, query, req.Name, req.Email, req.Role).StructScan(&user)
	if err != nil {
		r.logger.Error("Failed to create user",
			zap.String("email", req.Email),
			zap.Error(err))
		return nil, fmt.Errorf("create user: %w", err)
	}

	r.logger.Info("User created",
		zap.Int64("id", user.ID),
		zap.String("email", user.Email))

	return &user, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, name, email, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get user by ID",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

func (r *repository) List(ctx context.Context, filter *UserFilter) (*common.PaginatedResult[User], error) {
	filter.SetDefaults()

	var (
		conditions []string
		args       []interface{}
		argIndex   = 1
	)

	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIndex))
		args = append(args, "%"+filter.Name+"%")
		argIndex++
	}

	if filter.Role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, filter.Role)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		r.logger.Error("Failed to count users", zap.Error(err))
		return nil, fmt.Errorf("count users: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, name, email, role, created_at, updated_at
		FROM users
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	var users []User
	err = r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		r.logger.Error("Failed to list users", zap.Error(err))
		return nil, fmt.Errorf("list users: %w", err)
	}

	return &common.PaginatedResult[User]{
		Data:       users,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(filter.Limit))),
	}, nil
}

func (r *repository) Update(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error) {
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

	if req.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, *req.Email)
		argIndex++
	}

	if req.Role != nil {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, *req.Role)
		argIndex++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	setClauses = append(setClauses, "updated_at = now()")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = $%d
		RETURNING id, name, email, role, created_at, updated_at
	`, strings.Join(setClauses, ", "), argIndex)

	var user User
	err := r.db.QueryRowxContext(ctx, query, args...).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %d", id)
		}
		r.logger.Error("Failed to update user",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("update user: %w", err)
	}

	r.logger.Info("User updated", zap.Int64("id", user.ID))

	return &user, nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete user",
			zap.Int64("id", id),
			zap.Error(err))
		return fmt.Errorf("delete user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %d", id)
	}

	r.logger.Info("User deleted", zap.Int64("id", id))

	return nil
}
