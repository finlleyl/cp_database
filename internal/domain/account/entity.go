package account

import (
	"time"
)

// Account represents a trading account
type Account struct {
	ID              int64     `json:"id" db:"id"`
	UserID          int64     `json:"user_id" db:"user_id"`
	MTLogin         string    `json:"mt_login" db:"mt_login"`
	MTServerID      int64     `json:"mt_server_id" db:"mt_server_id"`
	SwapFreeStatus  bool      `json:"swap_free_status" db:"swap_free_status"`
	Balance         float64   `json:"balance" db:"balance"`
	Equity          float64   `json:"equity" db:"equity"`
	IsDeleted       bool      `json:"is_deleted" db:"is_deleted"`
	CommonCreatedAt time.Time `json:"common_created_at" db:"common_created_at"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// CreateAccountRequest represents the request to create a new account
type CreateAccountRequest struct {
	UserID          int64     `json:"user_id" binding:"required"`
	MTLogin         string    `json:"mt_login" binding:"required"`
	MTServerID      int64     `json:"mt_server_id" binding:"required"`
	SwapFreeStatus  bool      `json:"swap_free_status"`
	CommonCreatedAt time.Time `json:"common_created_at"`
}

// UpdateAccountRequest represents the request to update an account
type UpdateAccountRequest struct {
	MTLogin        *string  `json:"mt_login,omitempty"`
	MTServerID     *int64   `json:"mt_server_id,omitempty"`
	SwapFreeStatus *bool    `json:"swap_free_status,omitempty"`
	Balance        *float64 `json:"balance,omitempty"`
	Equity         *float64 `json:"equity,omitempty"`
}

// AccountFilter represents filter parameters for account search
type AccountFilter struct {
	UserID int64 `form:"user_id"`
}
