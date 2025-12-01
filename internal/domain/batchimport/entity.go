package batchimport

import (
	"encoding/json"
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

// ImportJob represents a batch import job
type ImportJob struct {
	ID               int64                  `json:"id" db:"id"`
	Type             ImportType             `json:"type" db:"type"`
	Status           common.ImportJobStatus `json:"status" db:"status"`
	FileName         string                 `json:"file_name" db:"file_name"`
	FileSize         int64                  `json:"file_size" db:"file_size"`
	TotalRecords     int                    `json:"total_records" db:"total_records"`
	ProcessedRecords int                    `json:"processed_records" db:"processed_records"`
	SuccessRecords   int                    `json:"success_records" db:"success_records"`
	FailedRecords    int                    `json:"failed_records" db:"failed_records"`
	Parameters       json.RawMessage        `json:"parameters" db:"parameters"`
	Summary          json.RawMessage        `json:"summary" db:"summary"`
	StartedAt        *time.Time             `json:"started_at,omitempty" db:"started_at"`
	CompletedAt      *time.Time             `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at"`
}

// ImportType represents the type of import
type ImportType string

const (
	ImportTypeTrades ImportType = "trades"
)

// ImportError represents an error that occurred during import
type ImportError struct {
	ID        int64     `json:"id" db:"id"`
	JobID     int64     `json:"job_id" db:"job_id"`
	RowNumber int       `json:"row_number" db:"row_number"`
	ErrorCode string    `json:"error_code" db:"error_code"`
	ErrorMsg  string    `json:"error_msg" db:"error_msg"`
	RawData   string    `json:"raw_data" db:"raw_data"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ImportJobSummary represents a summary of an import job
type ImportJobSummary struct {
	TotalRecords     int            `json:"total_records"`
	ProcessedRecords int            `json:"processed_records"`
	SuccessRecords   int            `json:"success_records"`
	FailedRecords    int            `json:"failed_records"`
	ErrorsByType     map[string]int `json:"errors_by_type"`
	Duration         time.Duration  `json:"duration"`
}

// ImportTradesRequest represents the request to import trades
type ImportTradesRequest struct {
	StrategyUUID string `form:"strategy_uuid" binding:"required"`
	AccountID    int64  `form:"account_id" binding:"required"`
	FileFormat   string `form:"file_format" binding:"required,oneof=csv json"`
}

// ImportTradesParameters represents parameters for trade import
type ImportTradesParameters struct {
	StrategyUUID string `json:"strategy_uuid"`
	AccountID    int64  `json:"account_id"`
	FileFormat   string `json:"file_format"`
}

// JobFilter represents filter parameters for import job search
type JobFilter struct {
	Type   ImportType             `form:"type"`
	Status common.ImportJobStatus `form:"status"`
	common.Pagination
}
