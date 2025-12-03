package httpserver

import (
	"net/http"

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
)

type RouteParams struct {
	UserHandler         *user.Handler
	AccountHandler      *account.Handler
	StrategyHandler     *strategy.Handler
	OfferHandler        *offer.Handler
	SubscriptionHandler *subscription.Handler
	TradeHandler        *trade.Handler
	StatisticsHandler   *statistics.Handler
	BatchImportHandler  *batchimport.Handler
	AuditHandler        *audit.Handler
}

func healthRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func pingRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func RegisterRoutes(r *gin.Engine, params RouteParams) {

	r.GET("/health", healthRoute)
	r.GET("/ping", pingRoute)

	v1 := r.Group("/api/v1")
	{

		user.RegisterRoutes(v1, params.UserHandler)
		account.RegisterRoutes(v1, params.AccountHandler)
		strategy.RegisterRoutes(v1, params.StrategyHandler)
		offer.RegisterRoutes(v1, params.OfferHandler)
		subscription.RegisterRoutes(v1, params.SubscriptionHandler)
		trade.RegisterRoutes(v1, params.TradeHandler)
		statistics.RegisterRoutes(v1, params.StatisticsHandler)
		batchimport.RegisterRoutes(v1, params.BatchImportHandler)
		audit.RegisterRoutes(v1, params.AuditHandler)
	}
}
