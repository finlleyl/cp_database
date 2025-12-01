package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/finlleyl/cp_database/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewDB(
	lc fx.Lifecycle,
	logger *zap.Logger,
	cfg *config.Config,
) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.POSTGRES_USER,
		cfg.POSTGRES_PASSWORD,
		"postgres",
		5432,
		cfg.POSTGRES_DB,
		"disable",
	)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx open: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Connecting to PostgreSQL")

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := db.PingContext(ctx); err != nil {
				logger.Error("Failed to ping PostgreSQL", zap.Error(err))
				return fmt.Errorf("ping postgres: %w", err)
			}

			logger.Info("PostgreSQL connected successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing PostgreSQL connection")
			if err := db.Close(); err != nil {
				logger.Error("Failed to close PostgreSQL connection", zap.Error(err))
				return err
			}
			return nil
		},
	})

	return db, nil
}

// Module provides database connection
// Domain repositories are provided by their respective modules
var Module = fx.Options(
	fx.Provide(
		NewDB,
	),
)
