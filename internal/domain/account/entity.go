package account

import (
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

type Account struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	AccountType string    `json:"account_type" db:"account_type"`
	Currency    string    `json:"currency" db:"currency"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateAccountRequest struct {
	UserID      int64  `json:"user_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	AccountType string `json:"account_type" binding:"required,oneof=master investor"`
	Currency    string `json:"currency" binding:"required,len=3"`
}

type UpdateAccountRequest struct {
	Name        *string `json:"name,omitempty"`
	AccountType *string `json:"account_type,omitempty"`
	Currency    *string `json:"currency,omitempty"`
}

type AccountFilter struct {
	UserID      int64  `form:"user_id"`
	AccountType string `form:"account_type"`
	common.Pagination
}

// AccountListResponse представляет пагинированный ответ со списком аккаунтов
type AccountListResponse struct {
	Data       []Account `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}
