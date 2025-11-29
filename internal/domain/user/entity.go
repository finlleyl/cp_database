package user

import (
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

// User represents a user in the copytrading platform
type User struct {
	ID        int64           `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	Email     string          `json:"email" db:"email"`
	Role      common.UserRole `json:"role" db:"role"`
	IsDeleted bool            `json:"is_deleted" db:"is_deleted"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Name  string          `json:"name" binding:"required"`
	Email string          `json:"email" binding:"required,email"`
	Role  common.UserRole `json:"role" binding:"required,oneof=master investor"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name  *string          `json:"name,omitempty"`
	Email *string          `json:"email,omitempty"`
	Role  *common.UserRole `json:"role,omitempty"`
}

// UserFilter represents filter parameters for user search
type UserFilter struct {
	Name string          `form:"name"`
	Role common.UserRole `form:"role"`
	common.Pagination
}
