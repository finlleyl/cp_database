package user

import (
	"context"
	"fmt"

	"github.com/finlleyl/cp_database/internal/domain/audit"
	"github.com/finlleyl/cp_database/internal/domain/common"
	"go.uber.org/zap"
)

// UseCase defines the interface for user business logic
type UseCase interface {
	Create(ctx context.Context, req *CreateUserRequest) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	List(ctx context.Context, filter *UserFilter) (*common.PaginatedResult[User], error)
	Update(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error)
	Delete(ctx context.Context, id int64) error
}

type useCase struct {
	repo      Repository
	auditRepo audit.Repository
	logger    *zap.Logger
}

// NewUseCase creates a new user use case
func NewUseCase(repo Repository, auditRepo audit.Repository, logger *zap.Logger) UseCase {
	return &useCase{repo: repo, auditRepo: auditRepo, logger: logger}
}

func (u *useCase) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
	// TODO: Implement user creation business logic
	// 1. Validate request
	// 2. Create user via repository
	// 3. Create audit log
	u.logger.Info("UseCase: Creating user", zap.String("name", req.Name))

	user, err := u.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeUser,
		EntityID:   fmt.Sprintf("%d", user.ID),
		Action:     audit.AuditActionCreate,
		NewValue:   user,
	})

	return user, nil
}

func (u *useCase) GetByID(ctx context.Context, id int64) (*User, error) {
	// TODO: Implement get user by ID business logic
	u.logger.Info("UseCase: Getting user by ID", zap.Int64("id", id))
	return u.repo.GetByID(ctx, id)
}

func (u *useCase) List(ctx context.Context, filter *UserFilter) (*common.PaginatedResult[User], error) {
	// TODO: Implement user listing business logic
	filter.SetDefaults()
	u.logger.Info("UseCase: Listing users", zap.Any("filter", filter))
	return u.repo.List(ctx, filter)
}

func (u *useCase) Update(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error) {
	// TODO: Implement user update business logic
	// 1. Get existing user
	// 2. Validate changes
	// 3. Update user
	// 4. Create audit log with old/new values
	u.logger.Info("UseCase: Updating user", zap.Int64("id", id))

	oldUser, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	user, err := u.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeUser,
		EntityID:   fmt.Sprintf("%d", id),
		Action:     audit.AuditActionUpdate,
		OldValue:   oldUser,
		NewValue:   user,
	})

	return user, nil
}

func (u *useCase) Delete(ctx context.Context, id int64) error {
	// TODO: Implement logical user deletion business logic
	// 1. Check for dependencies
	// 2. Mark user as deleted
	// 3. Create audit log
	u.logger.Info("UseCase: Deleting user", zap.Int64("id", id))

	oldUser, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeUser,
		EntityID:   fmt.Sprintf("%d", id),
		Action:     audit.AuditActionDelete,
		OldValue:   oldUser,
	})

	return nil
}
