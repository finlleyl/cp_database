package account

import (
	"context"
	"fmt"

	"github.com/finlleyl/cp_database/internal/domain/audit"
	"github.com/finlleyl/cp_database/internal/domain/common"
	"go.uber.org/zap"
)

// UseCase defines the interface for account business logic
type UseCase interface {
	Create(ctx context.Context, req *CreateAccountRequest) (*Account, error)
	GetByID(ctx context.Context, id int64) (*Account, error)
	List(ctx context.Context, filter *AccountFilter) (*common.PaginatedResult[Account], error)
	Update(ctx context.Context, id int64, req *UpdateAccountRequest) (*Account, error)
	Delete(ctx context.Context, id int64) error
}

type useCase struct {
	repo      Repository
	auditRepo audit.Repository
	logger    *zap.Logger
}

// NewUseCase creates a new account use case
func NewUseCase(repo Repository, auditRepo audit.Repository, logger *zap.Logger) UseCase {
	return &useCase{repo: repo, auditRepo: auditRepo, logger: logger}
}

func (u *useCase) Create(ctx context.Context, req *CreateAccountRequest) (*Account, error) {
	// TODO: Implement account creation business logic
	// 1. Validate user exists
	// 2. Check MT login uniqueness
	// 3. Create account
	// 4. Create audit log
	u.logger.Info("UseCase: Creating account",
		zap.Int64("user_id", req.UserID),
		zap.String("mt_login", req.MTLogin))

	account, err := u.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeAccount,
		EntityID:   account.ID,
		Action:     audit.AuditActionCreate,
		NewValue:   account,
	})

	return account, nil
}

func (u *useCase) GetByID(ctx context.Context, id int64) (*Account, error) {
	// TODO: Implement get account by ID business logic
	u.logger.Info("UseCase: Getting account by ID", zap.Int64("id", id))
	return u.repo.GetByID(ctx, id)
}

func (u *useCase) List(ctx context.Context, filter *AccountFilter) (*common.PaginatedResult[Account], error) {
	// TODO: Implement account listing business logic
	u.logger.Info("UseCase: Listing accounts", zap.Any("filter", filter))
	return u.repo.List(ctx, filter)
}

func (u *useCase) Update(ctx context.Context, id int64, req *UpdateAccountRequest) (*Account, error) {
	// TODO: Implement account update business logic
	// 1. Get existing account
	// 2. Validate changes
	// 3. Update account
	// 4. Create audit log
	u.logger.Info("UseCase: Updating account", zap.Int64("id", id))

	oldAccount, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}

	account, err := u.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("update account: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeAccount,
		EntityID:   id,
		Action:     audit.AuditActionUpdate,
		OldValue:   oldAccount,
		NewValue:   account,
	})

	return account, nil
}

func (u *useCase) Delete(ctx context.Context, id int64) error {
	// TODO: Implement account deletion business logic
	// 1. Check for dependent strategies/subscriptions
	// 2. Delete account if no dependencies
	// 3. Create audit log
	u.logger.Info("UseCase: Deleting account", zap.Int64("id", id))

	oldAccount, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete account: %w", err)
	}

	// Create audit log
	_, _ = u.auditRepo.Create(ctx, &audit.AuditCreateRequest{
		EntityType: audit.EntityTypeAccount,
		EntityID:   id,
		Action:     audit.AuditActionDelete,
		OldValue:   oldAccount,
	})

	return nil
}
