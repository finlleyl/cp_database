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
	"go.uber.org/zap"
)

type UseCase interface {

	CreateJob(ctx context.Context, req *CreateImportJobRequest) (*ImportJob, error)

	ImportTrades(ctx context.Context, req *ImportTradesRequest, file io.Reader, fileName string) (*ImportJob, error)

	GetJobByID(ctx context.Context, id int64) (*ImportJob, error)

	ListJobs(ctx context.Context, filter *JobFilter) (*common.PaginatedResult[ImportJob], error)

	GetJobErrors(ctx context.Context, jobID int64, filter *ErrorFilter) (*common.PaginatedResult[ImportJobError], error)

	GetJobSummary(ctx context.Context, jobID int64) (*ImportJobSummary, error)
}

type useCase struct {
	repo      Repository
	tradeRepo trade.Repository
	logger    *zap.Logger
}

func NewUseCase(repo Repository, tradeRepo trade.Repository, logger *zap.Logger) UseCase {
	return &useCase{repo: repo, tradeRepo: tradeRepo, logger: logger}
}

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

func (u *useCase) ImportTrades(ctx context.Context, req *ImportTradesRequest, file io.Reader, fileName string) (*ImportJob, error) {

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
		zap.Int64("strategy_id", req.StrategyID),
		zap.String("file_format", req.FileFormat))

	data, err := io.ReadAll(file)
	if err != nil {
		u.completeJobWithError(ctx, createdJob.ID, "Failed to read file: "+err.Error())
		return createdJob, nil
	}

	go u.processTradeImport(context.Background(), createdJob.ID, req, data)

	return createdJob, nil
}

func (u *useCase) processTradeImport(ctx context.Context, jobID int64, req *ImportTradesRequest, data []byte) {
	startTime := time.Now()

	var records []map[string]string
	var err error

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

	if err := u.repo.StartJob(ctx, jobID, totalRows); err != nil {
		u.logger.Error("Failed to start import job",
			zap.Int64("job_id", jobID),
			zap.Error(err))
		return
	}

	var processedRows, errorRows int
	var jobErrors []*ImportJobError

	for i, record := range records {
		rowNumber := i + 1

		createReq, err := u.mapRecordToTradeRequest(record, req.StrategyID, req.AccountID)
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

		if processedRows%100 == 0 {
			_ = u.repo.UpdateJobProgress(ctx, jobID, processedRows, errorRows)
		}
	}

	if len(jobErrors) > 0 {
		if err := u.repo.CreateErrorsBatch(ctx, jobErrors); err != nil {
			u.logger.Error("Failed to save import errors",
				zap.Int64("job_id", jobID),
				zap.Error(err))
		}
	}

	_ = u.repo.UpdateJobProgress(ctx, jobID, processedRows, errorRows)

	finalStatus := common.ImportJobStatusSuccess
	if errorRows > 0 && processedRows == 0 {
		finalStatus = common.ImportJobStatusFailed
	}

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

		var rawRecords []map[string]interface{}
		if err := json.Unmarshal(data, &rawRecords); err != nil {
			return nil, fmt.Errorf("parse JSON: %w", err)
		}

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

func (u *useCase) mapRecordToTradeRequest(record map[string]string, strategyID int64, accountID int64) (*trade.CreateTradeRequest, error) {

	symbol, ok := record["symbol"]
	if !ok || symbol == "" {
		return nil, fmt.Errorf("missing required field: symbol")
	}

	directionStr, ok := record["direction"]
	if !ok || directionStr == "" {

		directionStr, ok = record["type"]
		if !ok || directionStr == "" {
			return nil, fmt.Errorf("missing required field: direction")
		}
	}

	volumeStr, ok := record["volume_lots"]
	if !ok || volumeStr == "" {

		volumeStr, ok = record["volume"]
		if !ok || volumeStr == "" {
			return nil, fmt.Errorf("missing required field: volume_lots")
		}
	}

	volume, err := strconv.ParseFloat(volumeStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid volume_lots: %w", err)
	}

	openPriceStr, ok := record["open_price"]
	if !ok || openPriceStr == "" {
		return nil, fmt.Errorf("missing required field: open_price")
	}

	openPrice, err := strconv.ParseFloat(openPriceStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid open_price: %w", err)
	}

	openTimeStr, ok := record["open_time"]
	if !ok || openTimeStr == "" {
		return nil, fmt.Errorf("missing required field: open_time")
	}

	openTime, err := time.Parse(time.RFC3339, openTimeStr)
	if err != nil {

		openTime, err = time.Parse("2006-01-02 15:04:05", openTimeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid open_time format: %w", err)
		}
	}

	return &trade.CreateTradeRequest{
		StrategyID:      strategyID,
		MasterAccountID: accountID,
		Symbol:          symbol,
		Direction:       trade.TradeDirection(directionStr),
		VolumeLots:      volume,
		OpenPrice:       openPrice,
		OpenTime:        openTime,
	}, nil
}

func (u *useCase) completeJobWithError(ctx context.Context, jobID int64, errorMsg string) {

	jobError := &ImportJobError{
		JobID:        jobID,
		ErrorMessage: errorMsg,
	}
	_, _ = u.repo.CreateError(ctx, jobError)

	_ = u.repo.CompleteJob(ctx, jobID, common.ImportJobStatusFailed)

	u.logger.Error("Import job failed",
		zap.Int64("job_id", jobID),
		zap.String("error", errorMsg))
}

func (u *useCase) GetJobByID(ctx context.Context, id int64) (*ImportJob, error) {
	job, err := u.repo.GetJobByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get job by id: %w", err)
	}
	return job, nil
}

func (u *useCase) ListJobs(ctx context.Context, filter *JobFilter) (*common.PaginatedResult[ImportJob], error) {
	filter.SetDefaults()
	return u.repo.ListJobs(ctx, filter)
}

func (u *useCase) GetJobErrors(ctx context.Context, jobID int64, filter *ErrorFilter) (*common.PaginatedResult[ImportJobError], error) {
	filter.SetDefaults()
	return u.repo.GetJobErrors(ctx, jobID, filter)
}

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
