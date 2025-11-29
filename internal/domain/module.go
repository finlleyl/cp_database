package domain

import (
	"github.com/finlleyl/cp_database/internal/domain/account"
	"github.com/finlleyl/cp_database/internal/domain/audit"
	"github.com/finlleyl/cp_database/internal/domain/batchimport"
	"github.com/finlleyl/cp_database/internal/domain/offer"
	"github.com/finlleyl/cp_database/internal/domain/statistics"
	"github.com/finlleyl/cp_database/internal/domain/strategy"
	"github.com/finlleyl/cp_database/internal/domain/subscription"
	"github.com/finlleyl/cp_database/internal/domain/trade"
	"github.com/finlleyl/cp_database/internal/domain/user"
	"go.uber.org/fx"
)

// Module aggregates all domain modules
var Module = fx.Options(
	// Core domains
	user.Module,
	account.Module,

	// Business domains
	strategy.Module,
	offer.Module,
	subscription.Module,
	trade.Module,

	// Supporting domains
	statistics.Module,
	batchimport.Module,
	audit.Module,
)
