package batchimport

import (
	"encoding/json"
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

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

type ImportJobType string

const (
	ImportJobTypeTrades     ImportJobType = "trades"
	ImportJobTypeAccounts   ImportJobType = "accounts"
	ImportJobTypeStatistics ImportJobType = "statistics"
)

type ImportJobError struct {
	ID           int64           `json:"id" db:"id"`
	JobID        int64           `json:"job_id" db:"job_id"`
	RowNumber    *int            `json:"row_number,omitempty" db:"row_number"`
	RawData      json.RawMessage `json:"raw_data,omitempty" db:"raw_data"`
	ErrorMessage string          `json:"error_message" db:"error_message"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
}

type ImportJobSummary struct {
	TotalRows     int            `json:"total_rows"`
	ProcessedRows int            `json:"processed_rows"`
	ErrorRows     int            `json:"error_rows"`
	ErrorsByType  map[string]int `json:"errors_by_type,omitempty"`
	Duration      time.Duration  `json:"duration"`
}

type CreateImportJobRequest struct {
	Type     ImportJobType `json:"type" binding:"required,oneof=trades accounts statistics"`
	FileName string        `json:"file_name"`
}

type ImportTradesRequest struct {
	StrategyID int64  `form:"strategy_id" binding:"required"`
	AccountID  int64  `form:"account_id" binding:"required"`
	FileFormat string `form:"file_format" binding:"required,oneof=csv json"`
}

type ImportTradesParameters struct {
	StrategyID int64  `json:"strategy_id"`
	AccountID  int64  `json:"account_id"`
	FileFormat string `json:"file_format"`
}

type JobFilter struct {
	Type   ImportJobType          `form:"type"`
	Status common.ImportJobStatus `form:"status"`
	common.Pagination
}

type ErrorFilter struct {
	common.Pagination
}
