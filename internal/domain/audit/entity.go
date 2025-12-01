package audit

import (
	"encoding/json"
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID         int64           `json:"id" db:"id"`
	EntityType EntityType      `json:"entity_type" db:"entity_type"`
	EntityID   string          `json:"entity_id" db:"entity_id"`
	Action     AuditAction     `json:"action" db:"action"`
	UserID     *int64          `json:"user_id,omitempty" db:"user_id"`
	OldValue   json.RawMessage `json:"old_value,omitempty" db:"old_value"`
	NewValue   json.RawMessage `json:"new_value,omitempty" db:"new_value"`
	Changes    json.RawMessage `json:"changes,omitempty" db:"changes"`
	IPAddress  string          `json:"ip_address" db:"ip_address"`
	UserAgent  string          `json:"user_agent" db:"user_agent"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

// EntityType represents the type of entity being audited
type EntityType string

const (
	EntityTypeUser         EntityType = "users"
	EntityTypeAccount      EntityType = "accounts"
	EntityTypeStrategy     EntityType = "strategies"
	EntityTypeOffer        EntityType = "offers"
	EntityTypeSubscription EntityType = "subscriptions"
	EntityTypeTrade        EntityType = "trades"
)

// AuditAction represents the action being audited
type AuditAction string

const (
	AuditActionCreate       AuditAction = "create"
	AuditActionUpdate       AuditAction = "update"
	AuditActionDelete       AuditAction = "delete"
	AuditActionStatusChange AuditAction = "status_change"
)

// AuditFilter represents filter parameters for audit log search
type AuditFilter struct {
	Entity   EntityType  `form:"entity" binding:"omitempty,oneof=users accounts strategies offers subscriptions trades"`
	EntityID string      `form:"entity_id"`
	Action   AuditAction `form:"action" binding:"omitempty,oneof=create update delete status_change"`
	UserID   *int64      `form:"user_id"`
	common.TimeRange
	common.Pagination
}

// AuditCreateRequest is used internally to create audit log entries
type AuditCreateRequest struct {
	EntityType EntityType
	EntityID   int64
	Action     AuditAction
	UserID     *int64
	OldValue   interface{}
	NewValue   interface{}
	Changes    interface{}
	IPAddress  string
	UserAgent  string
}
