package batchimport

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Repository defines the interface for batch import data operations
type Repository interface {
	// Import Job operations
	CreateJob(ctx context.Context, job *ImportJob) (*ImportJob, error)
	GetJobByID(ctx context.Context, id int64) (*ImportJob, error)
	ListJobs(ctx context.Context, filter *JobFilter) (*common.PaginatedResult[ImportJob], error)
	UpdateJobStatus(ctx context.Context, id int64, status common.ImportJobStatus) error
	UpdateJobProgress(ctx context.Context, id int64, processedRows, errorRows int) error
	StartJob(ctx context.Context, id int64, totalRows int) error
	CompleteJob(ctx context.Context, id int64, status common.ImportJobStatus) error

	// Import Job Error operations
	CreateError(ctx context.Context, jobError *ImportJobError) (*ImportJobError, error)
	CreateErrorsBatch(ctx context.Context, jobErrors []*ImportJobError) error
	GetJobErrors(ctx context.Context, jobID int64, filter *ErrorFilter) (*common.PaginatedResult[ImportJobError], error)
	CountJobErrors(ctx context.Context, jobID int64) (int64, error)
}

type repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewRepository creates a new batch import repository
func NewRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &repository{db: db, logger: logger}
}

// CreateJob creates a new import job with pending status
func (r *repository) CreateJob(ctx context.Context, job *ImportJob) (*ImportJob, error) {
	query := `
		INSERT INTO import_jobs (type, status, file_name)
		VALUES ($1, $2, $3)
		RETURNING id, type, status, file_name, total_rows, processed_rows, error_rows, started_at, finished_at, created_at
	`

	var result ImportJob
	err := r.db.QueryRowxContext(ctx, query,
		job.Type,
		common.ImportJobStatusPending,
		job.FileName,
	).StructScan(&result)

	if err != nil {
		r.logger.Error("Failed to create import job",
			zap.String("type", string(job.Type)),
			zap.Error(err))
		return nil, fmt.Errorf("create import job: %w", err)
	}

	r.logger.Info("Import job created",
		zap.Int64("id", result.ID),
		zap.String("type", string(result.Type)))

	return &result, nil
}

// GetJobByID retrieves an import job by its ID
func (r *repository) GetJobByID(ctx context.Context, id int64) (*ImportJob, error) {
	query := `
		SELECT id, type, status, file_name, total_rows, processed_rows, error_rows, started_at, finished_at, created_at
		FROM import_jobs
		WHERE id = $1
	`

	var job ImportJob
	err := r.db.GetContext(ctx, &job, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get import job by ID",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("get import job by id: %w", err)
	}

	return &job, nil
}

// ListJobs retrieves a paginated list of import jobs with optional filters
func (r *repository) ListJobs(ctx context.Context, filter *JobFilter) (*common.PaginatedResult[ImportJob], error) {
	filter.SetDefaults()

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, filter.Type)
		argIndex++
	}

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filter.Status)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM import_jobs %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		r.logger.Error("Failed to count import jobs", zap.Error(err))
		return nil, fmt.Errorf("count import jobs: %w", err)
	}

	// Get data
	query := fmt.Sprintf(`
		SELECT id, type, status, file_name, total_rows, processed_rows, error_rows, started_at, finished_at, created_at
		FROM import_jobs
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filter.Limit, filter.Offset)

	var jobs []ImportJob
	err = r.db.SelectContext(ctx, &jobs, query, args...)
	if err != nil {
		r.logger.Error("Failed to list import jobs", zap.Error(err))
		return nil, fmt.Errorf("list import jobs: %w", err)
	}

	totalPages := int(total) / filter.Limit
	if int(total)%filter.Limit > 0 {
		totalPages++
	}

	return &common.PaginatedResult[ImportJob]{
		Data:       jobs,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateJobStatus updates the status of an import job
func (r *repository) UpdateJobStatus(ctx context.Context, id int64, status common.ImportJobStatus) error {
	query := `UPDATE import_jobs SET status = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		r.logger.Error("Failed to update import job status",
			zap.Int64("id", id),
			zap.String("status", string(status)),
			zap.Error(err))
		return fmt.Errorf("update import job status: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("import job not found: %d", id)
	}

	r.logger.Info("Import job status updated",
		zap.Int64("id", id),
		zap.String("status", string(status)))

	return nil
}

// UpdateJobProgress updates the progress counters of an import job
func (r *repository) UpdateJobProgress(ctx context.Context, id int64, processedRows, errorRows int) error {
	query := `
		UPDATE import_jobs 
		SET processed_rows = $1, error_rows = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, processedRows, errorRows, id)
	if err != nil {
		r.logger.Error("Failed to update import job progress",
			zap.Int64("id", id),
			zap.Int("processed_rows", processedRows),
			zap.Int("error_rows", errorRows),
			zap.Error(err))
		return fmt.Errorf("update import job progress: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("import job not found: %d", id)
	}

	return nil
}

// StartJob marks the job as running and sets started_at and total_rows
func (r *repository) StartJob(ctx context.Context, id int64, totalRows int) error {
	query := `
		UPDATE import_jobs 
		SET status = $1, started_at = $2, total_rows = $3
		WHERE id = $4
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, common.ImportJobStatusRunning, now, totalRows, id)
	if err != nil {
		r.logger.Error("Failed to start import job",
			zap.Int64("id", id),
			zap.Int("total_rows", totalRows),
			zap.Error(err))
		return fmt.Errorf("start import job: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("import job not found: %d", id)
	}

	r.logger.Info("Import job started",
		zap.Int64("id", id),
		zap.Int("total_rows", totalRows))

	return nil
}

// CompleteJob marks the job as completed (success or failed) and sets finished_at
func (r *repository) CompleteJob(ctx context.Context, id int64, status common.ImportJobStatus) error {
	query := `
		UPDATE import_jobs 
		SET status = $1, finished_at = $2
		WHERE id = $3
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, status, now, id)
	if err != nil {
		r.logger.Error("Failed to complete import job",
			zap.Int64("id", id),
			zap.String("status", string(status)),
			zap.Error(err))
		return fmt.Errorf("complete import job: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("import job not found: %d", id)
	}

	r.logger.Info("Import job completed",
		zap.Int64("id", id),
		zap.String("status", string(status)))

	return nil
}

// CreateError creates a single import error record
func (r *repository) CreateError(ctx context.Context, jobError *ImportJobError) (*ImportJobError, error) {
	query := `
		INSERT INTO import_job_errors (job_id, row_number, raw_data, error_message)
		VALUES ($1, $2, $3, $4)
		RETURNING id, job_id, row_number, raw_data, error_message, created_at
	`

	var result ImportJobError
	err := r.db.QueryRowxContext(ctx, query,
		jobError.JobID,
		jobError.RowNumber,
		jobError.RawData,
		jobError.ErrorMessage,
	).StructScan(&result)

	if err != nil {
		r.logger.Error("Failed to create import error",
			zap.Int64("job_id", jobError.JobID),
			zap.Error(err))
		return nil, fmt.Errorf("create import error: %w", err)
	}

	return &result, nil
}

// CreateErrorsBatch creates multiple import error records in a single transaction
func (r *repository) CreateErrorsBatch(ctx context.Context, jobErrors []*ImportJobError) error {
	if len(jobErrors) == 0 {
		return nil
	}

	query := `
		INSERT INTO import_job_errors (job_id, row_number, raw_data, error_message)
		VALUES ($1, $2, $3, $4)
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, jobError := range jobErrors {
		_, err = stmt.ExecContext(ctx,
			jobError.JobID,
			jobError.RowNumber,
			jobError.RawData,
			jobError.ErrorMessage,
		)
		if err != nil {
			r.logger.Error("Failed to insert import error in batch",
				zap.Int64("job_id", jobError.JobID),
				zap.Error(err))
			return fmt.Errorf("insert import error: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	r.logger.Info("Import errors batch created",
		zap.Int("count", len(jobErrors)))

	return nil
}

// GetJobErrors retrieves paginated errors for a specific job
func (r *repository) GetJobErrors(ctx context.Context, jobID int64, filter *ErrorFilter) (*common.PaginatedResult[ImportJobError], error) {
	filter.SetDefaults()

	// Count total
	countQuery := `SELECT COUNT(*) FROM import_job_errors WHERE job_id = $1`
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, jobID)
	if err != nil {
		r.logger.Error("Failed to count import job errors",
			zap.Int64("job_id", jobID),
			zap.Error(err))
		return nil, fmt.Errorf("count import job errors: %w", err)
	}

	// Get data
	query := `
		SELECT id, job_id, row_number, raw_data, error_message, created_at
		FROM import_job_errors
		WHERE job_id = $1
		ORDER BY row_number ASC, created_at ASC
		LIMIT $2 OFFSET $3
	`

	var errors []ImportJobError
	err = r.db.SelectContext(ctx, &errors, query, jobID, filter.Limit, filter.Offset)
	if err != nil {
		r.logger.Error("Failed to get import job errors",
			zap.Int64("job_id", jobID),
			zap.Error(err))
		return nil, fmt.Errorf("get import job errors: %w", err)
	}

	totalPages := int(total) / filter.Limit
	if int(total)%filter.Limit > 0 {
		totalPages++
	}

	return &common.PaginatedResult[ImportJobError]{
		Data:       errors,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// CountJobErrors returns the total number of errors for a job
func (r *repository) CountJobErrors(ctx context.Context, jobID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM import_job_errors WHERE job_id = $1`

	var count int64
	err := r.db.GetContext(ctx, &count, query, jobID)
	if err != nil {
		r.logger.Error("Failed to count import job errors",
			zap.Int64("job_id", jobID),
			zap.Error(err))
		return 0, fmt.Errorf("count import job errors: %w", err)
	}

	return count, nil
}
