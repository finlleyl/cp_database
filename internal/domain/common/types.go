package common

import (
	"time"

	"github.com/google/uuid"
)

type Pagination struct {
	Page   int `json:"page" form:"page"`
	Limit  int `json:"limit" form:"limit"`
	Offset int `json:"-"`
}

func (p *Pagination) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	p.Offset = (p.Page - 1) * p.Limit
}

type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

type StrategyStatus string

const (
	StrategyStatusPreparing StrategyStatus = "preparing"
	StrategyStatusActive    StrategyStatus = "active"
	StrategyStatusArchived  StrategyStatus = "archived"
	StrategyStatusDeleted   StrategyStatus = "deleted"
)

type OfferStatus string

const (
	OfferStatusActive   OfferStatus = "active"
	OfferStatusArchived OfferStatus = "archived"
	OfferStatusDeleted  OfferStatus = "deleted"
)

type SubscriptionStatus string

const (
	SubscriptionStatusPreparing SubscriptionStatus = "preparing"
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusArchived  SubscriptionStatus = "archived"
	SubscriptionStatusSuspended SubscriptionStatus = "suspended"
	SubscriptionStatusDeleted   SubscriptionStatus = "deleted"
)

type ImportJobStatus string

const (
	ImportJobStatusPending ImportJobStatus = "pending"
	ImportJobStatusRunning ImportJobStatus = "running"
	ImportJobStatusSuccess ImportJobStatus = "success"
	ImportJobStatusFailed  ImportJobStatus = "failed"
)

type AuditOperation string

const (
	AuditOperationInsert AuditOperation = "insert"
	AuditOperationUpdate AuditOperation = "update"
	AuditOperationDelete AuditOperation = "delete"
)

type UserRole string

const (
	UserRoleMaster   UserRole = "master"
	UserRoleInvestor UserRole = "investor"
)

type FeeInterval string

const (
	FeeIntervalDaily   FeeInterval = "daily"
	FeeIntervalWeekly  FeeInterval = "weekly"
	FeeIntervalMonthly FeeInterval = "monthly"
)

func NewUUID() uuid.UUID {
	return uuid.New()
}

type TimeRange struct {
	From time.Time `json:"from" form:"from"`
	To   time.Time `json:"to" form:"to"`
}
