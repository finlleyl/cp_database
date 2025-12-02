package batchimport

import (
	"encoding/json"
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

// ImportJob represents a batch import job (matches import_jobs table)
type ImportJob struct {
	ID            int64                  `json:"id" db:"id"`
	Type          ImportJobType          `json:"type" db:"type"`
	Status        common.ImportJobStatus `json:"status" db:"status"`
	FileName      *string                `json:"file_name,omitempty" db:"file_name"`
	TotalRows     int                    `json:"total_rows" db:"total_rows"`
	ProcessedRows int                    `json:"processed_rows" db:"processed_rows"`
	ErrorRows     int                    `json:"error_rows" db:"error_rows"`
	StartedAt     *time.Time             `json:"started_at,omitempty" db:"started_at"`
	FinishedAt    *time.Time             `json:"finished_at,omitempty" db:"finished_at"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
}

// ImportJobType represents the type of import (matches import_job_type enum)
type ImportJobType string

const (
	ImportJobTypeTrades     ImportJobType = "trades"
	ImportJobTypeAccounts   ImportJobType = "accounts"
	ImportJobTypeStatistics ImportJobType = "statistics"
)

// ImportJobError represents an error that occurred during import (matches import_job_errors table)
type ImportJobError struct {
	ID           int64           `json:"id" db:"id"`
	JobID        int64           `json:"job_id" db:"job_id"`
	RowNumber    *int            `json:"row_number,omitempty" db:"row_number"`
	RawData      json.RawMessage `json:"raw_data,omitempty" db:"raw_data"`
	ErrorMessage string          `json:"error_message" db:"error_message"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
}

// ImportJobSummary represents a summary of an import job (computed, not stored)
type ImportJobSummary struct {
	TotalRows     int            `json:"total_rows"`
	ProcessedRows int            `json:"processed_rows"`
	ErrorRows     int            `json:"error_rows"`
	ErrorsByType  map[string]int `json:"errors_by_type,omitempty"`
	Duration      time.Duration  `json:"duration"`
}

// CreateImportJobRequest represents the request to create an import job
type CreateImportJobRequest struct {
	Type     ImportJobType `json:"type" binding:"required,oneof=trades accounts statistics"`
	FileName string        `json:"file_name"`
}

// ImportTradesRequest represents the request to import trades
type ImportTradesRequest struct {
	StrategyUUID string `form:"strategy_uuid" binding:"required"`
	AccountID    int64  `form:"account_id" binding:"required"`
	FileFormat   string `form:"file_format" binding:"required,oneof=csv json"`
}

// ImportTradesParameters represents parameters for trade import (stored as JSON metadata)
type ImportTradesParameters struct {
	StrategyUUID string `json:"strategy_uuid"`
	AccountID    int64  `json:"account_id"`
	FileFormat   string `json:"file_format"`
}

// JobFilter represents filter parameters for import job search
type JobFilter struct {
	Type   ImportJobType          `form:"type"`
	Status common.ImportJobStatus `form:"status"`
	common.Pagination
}

// ErrorFilter represents filter parameters for import job error search
type ErrorFilter struct {
	common.Pagination
}
