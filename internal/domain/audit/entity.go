package audit

import (
	"encoding/json"
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

type AuditLog struct {
	ID         int64                 `json:"id" db:"id"`
	EntityName string                `json:"entity_name" db:"entity_name"`
	EntityPK   string                `json:"entity_pk" db:"entity_pk"`
	Operation  common.AuditOperation `json:"operation" db:"operation"`
	ChangedBy  *int64                `json:"changed_by,omitempty" db:"changed_by"`
	ChangedAt  time.Time             `json:"changed_at" db:"changed_at"`
	OldRow     json.RawMessage       `json:"old_row,omitempty" db:"old_row" swaggertype:"object"`
	NewRow     json.RawMessage       `json:"new_row,omitempty" db:"new_row" swaggertype:"object"`
}

const (
	EntityNameUsers         = "users"
	EntityNameAccounts      = "accounts"
	EntityNameStrategies    = "strategies"
	EntityNameOffers        = "offers"
	EntityNameSubscriptions = "subscriptions"
	EntityNameTrades        = "trades"
)

type EntityType string

const (
	EntityTypeUser         EntityType = "users"
	EntityTypeAccount      EntityType = "accounts"
	EntityTypeStrategy     EntityType = "strategies"
	EntityTypeOffer        EntityType = "offers"
	EntityTypeSubscription EntityType = "subscriptions"
	EntityTypeTrade        EntityType = "trades"
)

type AuditAction string

const (
	AuditActionCreate       AuditAction = "insert"
	AuditActionUpdate       AuditAction = "update"
	AuditActionDelete       AuditAction = "delete"
	AuditActionStatusChange AuditAction = "update"
)

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

type AuditFilter struct {
	EntityName string                `form:"entity_name" binding:"omitempty,oneof=users accounts strategies offers subscriptions trades"`
	EntityPK   string                `form:"entity_pk"`
	Operation  common.AuditOperation `form:"operation" binding:"omitempty,oneof=insert update delete"`
	ChangedBy  *int64                `form:"changed_by"`
	common.TimeRange
	common.Pagination
}

type AuditStats struct {
	EntityName   string `json:"entity_name" db:"entity_name"`
	Operation    string `json:"operation" db:"operation"`
	TotalChanges int64  `json:"total_changes" db:"total_changes"`
}

type AuditStatsFilter struct {
	EntityName string `form:"entity_name"`
	common.TimeRange
}

// AuditListResponse представляет пагинированный ответ со списком аудит-логов
type AuditListResponse struct {
	Data       []AuditLog `json:"data"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int        `json:"total_pages"`
}
