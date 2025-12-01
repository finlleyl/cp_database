package batchimport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/finlleyl/cp_database/internal/domain/trade"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UseCase defines the interface for batch import business logic
type UseCase interface {
	ImportTrades(ctx context.Context, req *ImportTradesRequest, file io.Reader, fileName string, fileSize int64) (*ImportJob, error)
	GetJobByID(ctx context.Context, id int64) (*ImportJob, error)
	ListJobs(ctx context.Context, filter *JobFilter) (*common.PaginatedResult[ImportJob], error)
	GetJobErrors(ctx context.Context, jobID int64) ([]*ImportError, error)
}

type useCase struct {
	repo      Repository
	tradeRepo trade.Repository
	logger    *zap.Logger
}

// NewUseCase creates a new batch import use case
func NewUseCase(repo Repository, tradeRepo trade.Repository, logger *zap.Logger) UseCase {
	return &useCase{repo: repo, tradeRepo: tradeRepo, logger: logger}
}

func (u *useCase) ImportTrades(ctx context.Context, req *ImportTradesRequest, file io.Reader, fileName string, fileSize int64) (*ImportJob, error) {
	// TODO: Implement trade import business logic
	// Business flow:
	// 1. Create import job with status = pending
	// 2. Parse file (CSV or JSON)
	// 3. Update job status to processing
	// 4. For each record:
	//    a. Validate record
	//    b. Create trade
	//    c. Update job progress
	//    d. Record errors if any
	// 5. Complete job with summary
	u.logger.Info("UseCase: Importing trades",
		zap.String("strategy_uuid", req.StrategyUUID),
		zap.String("file_name", fileName),
		zap.String("file_format", req.FileFormat))

	params, _ := json.Marshal(&ImportTradesParameters{
		StrategyUUID: req.StrategyUUID,
		AccountID:    req.AccountID,
		FileFormat:   req.FileFormat,
	})

	job := &ImportJob{
		Type:       ImportTypeTrades,
		Status:     common.ImportJobStatusPending,
		FileName:   fileName,
		FileSize:   fileSize,
		Parameters: params,
	}

	job, err := u.repo.CreateJob(ctx, job)
	if err != nil {
		return nil, fmt.Errorf("create import job: %w", err)
	}

	// Start processing in background
	go u.processTradeImport(context.Background(), job.ID, req, file)

	return job, nil
}

func (u *useCase) processTradeImport(ctx context.Context, jobID int64, req *ImportTradesRequest, file io.Reader) {
	// TODO: Implement actual trade import processing
	// This is a placeholder for the background processing logic
	u.logger.Info("Processing trade import", zap.Int64("job_id", jobID))

	// Update status to processing
	if err := u.repo.UpdateJobStatus(ctx, jobID, common.ImportJobStatusProcessing); err != nil {
		u.logger.Error("Failed to update job status", zap.Error(err))
		return
	}

	startTime := time.Now()
	var processed, success, failed int

	// TODO: Parse file based on format (CSV or JSON)
	// TODO: For each record, create trade and track progress

	strategyUUID, _ := uuid.Parse(req.StrategyUUID)

	// Placeholder: simulate processing
	_ = strategyUUID
	processed = 0
	success = 0
	failed = 0

	// Complete job with summary
	summary := &ImportJobSummary{
		TotalRecords:     processed,
		ProcessedRecords: processed,
		SuccessRecords:   success,
		FailedRecords:    failed,
		ErrorsByType:     make(map[string]int),
		Duration:         time.Since(startTime),
	}

	if err := u.repo.CompleteJob(ctx, jobID, summary); err != nil {
		u.logger.Error("Failed to complete job", zap.Error(err))
	}

	u.logger.Info("Trade import completed",
		zap.Int64("job_id", jobID),
		zap.Int("processed", processed),
		zap.Int("success", success),
		zap.Int("failed", failed))
}

func (u *useCase) GetJobByID(ctx context.Context, id int64) (*ImportJob, error) {
	// TODO: Implement get import job by ID business logic
	u.logger.Info("UseCase: Getting import job by ID", zap.Int64("id", id))
	return u.repo.GetJobByID(ctx, id)
}

func (u *useCase) ListJobs(ctx context.Context, filter *JobFilter) (*common.PaginatedResult[ImportJob], error) {
	// TODO: Implement import job listing business logic
	filter.SetDefaults()
	u.logger.Info("UseCase: Listing import jobs", zap.Any("filter", filter))
	return u.repo.ListJobs(ctx, filter)
}

func (u *useCase) GetJobErrors(ctx context.Context, jobID int64) ([]*ImportError, error) {
	// TODO: Implement get import job errors business logic
	u.logger.Info("UseCase: Getting import job errors", zap.Int64("job_id", jobID))
	return u.repo.GetJobErrors(ctx, jobID)
}
