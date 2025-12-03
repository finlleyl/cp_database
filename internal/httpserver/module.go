package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/finlleyl/cp_database/internal/config"
	"github.com/finlleyl/cp_database/internal/domain/account"
	"github.com/finlleyl/cp_database/internal/domain/audit"
	"github.com/finlleyl/cp_database/internal/domain/batchimport"
	"github.com/finlleyl/cp_database/internal/domain/offer"
	"github.com/finlleyl/cp_database/internal/domain/statistics"
	"github.com/finlleyl/cp_database/internal/domain/strategy"
	"github.com/finlleyl/cp_database/internal/domain/subscription"
	"github.com/finlleyl/cp_database/internal/domain/trade"
	"github.com/finlleyl/cp_database/internal/domain/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger) *gin.Engine {
	r := gin.New()

	r.Use(
		gin.Recovery(),
		gin.Logger(),
		GinLogger(logger),
		RequestID(),
		CORSMiddleware(),
	)

	return r
}

func RegisterAllRoutes(
	r *gin.Engine,
	userHandler *user.Handler,
	accountHandler *account.Handler,
	strategyHandler *strategy.Handler,
	offerHandler *offer.Handler,
	subscriptionHandler *subscription.Handler,
	tradeHandler *trade.Handler,
	statisticsHandler *statistics.Handler,
	batchImportHandler *batchimport.Handler,
	auditHandler *audit.Handler,
) {
	params := RouteParams{
		UserHandler:         userHandler,
		AccountHandler:      accountHandler,
		StrategyHandler:     strategyHandler,
		OfferHandler:        offerHandler,
		SubscriptionHandler: subscriptionHandler,
		TradeHandler:        tradeHandler,
		StatisticsHandler:   statisticsHandler,
		BatchImportHandler:  batchImportHandler,
		AuditHandler:        auditHandler,
	}
	RegisterRoutes(r, params)
}

func NewHTTPServer(
	lc fx.Lifecycle,
	cfg *config.Config,
	router *gin.Engine,
	logger *zap.Logger,
) *http.Server {
	addr := fmt.Sprintf(":%s", cfg.HTTPPort)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("HTTP server starting", zap.String("addr", addr))

			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("HTTP server shutting down")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

var Module = fx.Options(
	fx.Provide(
		NewRouter,
		NewHTTPServer,
	),
	fx.Invoke(
		RegisterAllRoutes,
	),
)
