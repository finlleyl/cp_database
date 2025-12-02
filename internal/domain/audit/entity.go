package audit

import (
	"encoding/json"
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

// AuditLog represents an audit log entry (matches audit_log table)
type AuditLog struct {
	ID         int64                 `json:"id" db:"id"`
	EntityName string                `json:"entity_name" db:"entity_name"`
	EntityPK   string                `json:"entity_pk" db:"entity_pk"`
	Operation  common.AuditOperation `json:"operation" db:"operation"`
	ChangedBy  *int64                `json:"changed_by,omitempty" db:"changed_by"`
	ChangedAt  time.Time             `json:"changed_at" db:"changed_at"`
	OldRow     json.RawMessage       `json:"old_row,omitempty" db:"old_row"`
	NewRow     json.RawMessage       `json:"new_row,omitempty" db:"new_row"`
}

// Auditable entity names (matches tables with audit triggers)
const (
	EntityNameUsers         = "users"
	EntityNameAccounts      = "accounts"
	EntityNameStrategies    = "strategies"
	EntityNameOffers        = "offers"
	EntityNameSubscriptions = "subscriptions"
	EntityNameTrades        = "trades"
)

// EntityType represents the type of entity being audited (for backward compatibility)
type EntityType string

const (
	EntityTypeUser         EntityType = "users"
	EntityTypeAccount      EntityType = "accounts"
	EntityTypeStrategy     EntityType = "strategies"
	EntityTypeOffer        EntityType = "offers"
	EntityTypeSubscription EntityType = "subscriptions"
	EntityTypeTrade        EntityType = "trades"
)

// AuditAction represents the action being audited (for backward compatibility)
type AuditAction string

const (
	AuditActionCreate       AuditAction = "insert"
	AuditActionUpdate       AuditAction = "update"
	AuditActionDelete       AuditAction = "delete"
	AuditActionStatusChange AuditAction = "update" // Maps to update in DB
)

// AuditCreateRequest is used internally to create audit log entries (for backward compatibility)
// Note: With database triggers, manual audit creation is optional as triggers handle it automatically
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

// AuditFilter represents filter parameters for audit log search
type AuditFilter struct {
	EntityName string                `form:"entity_name" binding:"omitempty,oneof=users accounts strategies offers subscriptions trades"`
	EntityPK   string                `form:"entity_pk"`
	Operation  common.AuditOperation `form:"operation" binding:"omitempty,oneof=insert update delete"`
	ChangedBy  *int64                `form:"changed_by"`
	common.TimeRange
	common.Pagination
}

// AuditStats represents statistics for audit log
type AuditStats struct {
	EntityName   string `json:"entity_name" db:"entity_name"`
	Operation    string `json:"operation" db:"operation"`
	TotalChanges int64  `json:"total_changes" db:"total_changes"`
}

// AuditStatsFilter represents filter parameters for audit statistics
type AuditStatsFilter struct {
	EntityName string `form:"entity_name"`
	common.TimeRange
}
