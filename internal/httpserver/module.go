package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/finlleyl/cp_database/internal/config"
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

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", healthRoute)
	r.GET("/ping", pingRoute)
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
		RegisterRoutes,
	),
)
