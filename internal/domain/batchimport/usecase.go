package batchimport

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
	"github.com/finlleyl/cp_database/internal/domain/trade"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UseCase defines the interface for batch import business logic
type UseCase interface {
	// CreateJob creates a new import job
	CreateJob(ctx context.Context, req *CreateImportJobRequest) (*ImportJob, error)

	// ImportTrades starts a trade import process
	ImportTrades(ctx context.Context, req *ImportTradesRequest, file io.Reader, fileName string) (*ImportJob, error)

	// GetJobByID retrieves an import job by ID
	GetJobByID(ctx context.Context, id int64) (*ImportJob, error)

	// ListJobs retrieves a paginated list of import jobs
	ListJobs(ctx context.Context, filter *JobFilter) (*common.PaginatedResult[ImportJob], error)

	// GetJobErrors retrieves errors for a specific job
	GetJobErrors(ctx context.Context, jobID int64, filter *ErrorFilter) (*common.PaginatedResult[ImportJobError], error)

	// GetJobSummary retrieves summary statistics for a job
	GetJobSummary(ctx context.Context, jobID int64) (*ImportJobSummary, error)
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

// CreateJob creates a new import job with pending status
func (u *useCase) CreateJob(ctx context.Context, req *CreateImportJobRequest) (*ImportJob, error) {
	job := &ImportJob{
		Type:     req.Type,
		FileName: &req.FileName,
	}

	createdJob, err := u.repo.CreateJob(ctx, job)
	if err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}

	u.logger.Info("Import job created",
		zap.Int64("id", createdJob.ID),
		zap.String("type", string(createdJob.Type)))

	return createdJob, nil
}

// ImportTrades handles the trade import workflow
func (u *useCase) ImportTrades(ctx context.Context, req *ImportTradesRequest, file io.Reader, fileName string) (*ImportJob, error) {
	// Create the import job
	job := &ImportJob{
		Type:     ImportJobTypeTrades,
		FileName: &fileName,
	}

	createdJob, err := u.repo.CreateJob(ctx, job)
	if err != nil {
		return nil, fmt.Errorf("create import job: %w", err)
	}

	u.logger.Info("Trade import job created",
		zap.Int64("job_id", createdJob.ID),
		zap.String("strategy_uuid", req.StrategyUUID),
		zap.String("file_format", req.FileFormat))

	// Read file data before starting goroutine (file may be closed after handler returns)
	data, err := io.ReadAll(file)
	if err != nil {
		u.completeJobWithError(ctx, createdJob.ID, "Failed to read file: "+err.Error())
		return createdJob, nil
	}

	// Process import in background
	go u.processTradeImport(context.Background(), createdJob.ID, req, data)

	return createdJob, nil
}

// processTradeImport handles the actual trade import processing
func (u *useCase) processTradeImport(ctx context.Context, jobID int64, req *ImportTradesRequest, data []byte) {
	startTime := time.Now()

	var records []map[string]string
	var err error

	// Parse based on format
	switch req.FileFormat {
	case "csv":
		records, err = u.parseCSV(data)
	case "json":
		records, err = u.parseJSON(data)
	default:
		u.completeJobWithError(ctx, jobID, "Unsupported file format: "+req.FileFormat)
		return
	}

	if err != nil {
		u.completeJobWithError(ctx, jobID, "Failed to parse file: "+err.Error())
		return
	}

	totalRows := len(records)

	// Start the job with total row count
	if err := u.repo.StartJob(ctx, jobID, totalRows); err != nil {
		u.logger.Error("Failed to start import job",
			zap.Int64("job_id", jobID),
			zap.Error(err))
		return
	}

	strategyUUID, err := uuid.Parse(req.StrategyUUID)
	if err != nil {
		u.completeJobWithError(ctx, jobID, "Invalid strategy UUID: "+err.Error())
		return
	}

	var processedRows, errorRows int
	var jobErrors []*ImportJobError

	// Process each record
	for i, record := range records {
		rowNumber := i + 1

		createReq, err := u.mapRecordToTradeRequest(record, strategyUUID, req.AccountID)
		if err != nil {
			errorRows++
			rawData, _ := json.Marshal(record)
			jobErrors = append(jobErrors, &ImportJobError{
				JobID:        jobID,
				RowNumber:    &rowNumber,
				RawData:      rawData,
				ErrorMessage: err.Error(),
			})
			continue
		}

		// Create the trade
		_, err = u.tradeRepo.Create(ctx, createReq)
		if err != nil {
			errorRows++
			rawData, _ := json.Marshal(record)
			jobErrors = append(jobErrors, &ImportJobError{
				JobID:        jobID,
				RowNumber:    &rowNumber,
				RawData:      rawData,
				ErrorMessage: "Failed to create trade: " + err.Error(),
			})
			continue
		}

		processedRows++

		// Update progress every 100 records
		if processedRows%100 == 0 {
			_ = u.repo.UpdateJobProgress(ctx, jobID, processedRows, errorRows)
		}
	}

	// Save all errors in batch
	if len(jobErrors) > 0 {
		if err := u.repo.CreateErrorsBatch(ctx, jobErrors); err != nil {
			u.logger.Error("Failed to save import errors",
				zap.Int64("job_id", jobID),
				zap.Error(err))
		}
	}

	// Final progress update
	_ = u.repo.UpdateJobProgress(ctx, jobID, processedRows, errorRows)

	// Determine final status
	finalStatus := common.ImportJobStatusSuccess
	if errorRows > 0 && processedRows == 0 {
		finalStatus = common.ImportJobStatusFailed
	}

	// Complete the job
	if err := u.repo.CompleteJob(ctx, jobID, finalStatus); err != nil {
		u.logger.Error("Failed to complete import job",
			zap.Int64("job_id", jobID),
			zap.Error(err))
	}

	u.logger.Info("Trade import completed",
		zap.Int64("job_id", jobID),
		zap.Int("total_rows", totalRows),
		zap.Int("processed", processedRows),
		zap.Int("errors", errorRows),
		zap.Duration("duration", time.Since(startTime)))
}

func (u *useCase) parseCSV(data []byte) ([]map[string]string, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read CSV header: %w", err)
	}

	var records []map[string]string
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read CSV row: %w", err)
		}

		record := make(map[string]string)
		for i, value := range row {
			if i < len(header) {
				record[header[i]] = value
			}
		}
		records = append(records, record)
	}

	return records, nil
}

func (u *useCase) parseJSON(data []byte) ([]map[string]string, error) {
	var records []map[string]string
	if err := json.Unmarshal(data, &records); err != nil {
		// Try parsing as array of objects with any type
		var rawRecords []map[string]interface{}
		if err := json.Unmarshal(data, &rawRecords); err != nil {
			return nil, fmt.Errorf("parse JSON: %w", err)
		}

		// Convert to string map
		for _, raw := range rawRecords {
			record := make(map[string]string)
			for k, v := range raw {
				record[k] = fmt.Sprintf("%v", v)
			}
			records = append(records, record)
		}
	}

	return records, nil
}

func (u *useCase) mapRecordToTradeRequest(record map[string]string, strategyUUID uuid.UUID, accountID int64) (*trade.CreateTradeRequest, error) {
	// Parse required fields
	symbol, ok := record["symbol"]
	if !ok || symbol == "" {
		return nil, fmt.Errorf("missing required field: symbol")
	}

	typeStr, ok := record["type"]
	if !ok || typeStr == "" {
		return nil, fmt.Errorf("missing required field: type")
	}

	volumeStr, ok := record["volume"]
	if !ok || volumeStr == "" {
		return nil, fmt.Errorf("missing required field: volume")
	}

	volume, err := strconv.ParseFloat(volumeStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid volume: %w", err)
	}

	openPriceStr, ok := record["open_price"]
	if !ok || openPriceStr == "" {
		return nil, fmt.Errorf("missing required field: open_price")
	}

	openPrice, err := strconv.ParseFloat(openPriceStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid open_price: %w", err)
	}

	// Parse optional fields
	var stopLoss *float64
	if sl, ok := record["stop_loss"]; ok && sl != "" {
		parsed, err := strconv.ParseFloat(sl, 64)
		if err == nil {
			stopLoss = &parsed
		}
	}

	var takeProfit *float64
	if tp, ok := record["take_profit"]; ok && tp != "" {
		parsed, err := strconv.ParseFloat(tp, 64)
		if err == nil {
			takeProfit = &parsed
		}
	}

	return &trade.CreateTradeRequest{
		StrategyUUID: strategyUUID,
		AccountID:    accountID,
		Symbol:       symbol,
		Type:         trade.TradeType(typeStr),
		Volume:       volume,
		OpenPrice:    openPrice,
		StopLoss:     stopLoss,
		TakeProfit:   takeProfit,
	}, nil
}

func (u *useCase) completeJobWithError(ctx context.Context, jobID int64, errorMsg string) {
	// Create error record
	jobError := &ImportJobError{
		JobID:        jobID,
		ErrorMessage: errorMsg,
	}
	_, _ = u.repo.CreateError(ctx, jobError)

	// Complete job as failed
	_ = u.repo.CompleteJob(ctx, jobID, common.ImportJobStatusFailed)

	u.logger.Error("Import job failed",
		zap.Int64("job_id", jobID),
		zap.String("error", errorMsg))
}

// GetJobByID retrieves an import job by its ID
func (u *useCase) GetJobByID(ctx context.Context, id int64) (*ImportJob, error) {
	job, err := u.repo.GetJobByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get job by id: %w", err)
	}
	return job, nil
}

// ListJobs retrieves a paginated list of import jobs
func (u *useCase) ListJobs(ctx context.Context, filter *JobFilter) (*common.PaginatedResult[ImportJob], error) {
	filter.SetDefaults()
	return u.repo.ListJobs(ctx, filter)
}

// GetJobErrors retrieves errors for a specific job
func (u *useCase) GetJobErrors(ctx context.Context, jobID int64, filter *ErrorFilter) (*common.PaginatedResult[ImportJobError], error) {
	filter.SetDefaults()
	return u.repo.GetJobErrors(ctx, jobID, filter)
}

// GetJobSummary retrieves summary statistics for a completed job
func (u *useCase) GetJobSummary(ctx context.Context, jobID int64) (*ImportJobSummary, error) {
	job, err := u.repo.GetJobByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("get job: %w", err)
	}
	if job == nil {
		return nil, fmt.Errorf("job not found: %d", jobID)
	}

	var duration time.Duration
	if job.StartedAt != nil && job.FinishedAt != nil {
		duration = job.FinishedAt.Sub(*job.StartedAt)
	}

	return &ImportJobSummary{
		TotalRows:     job.TotalRows,
		ProcessedRows: job.ProcessedRows,
		ErrorRows:     job.ErrorRows,
		Duration:      duration,
	}, nil
}
