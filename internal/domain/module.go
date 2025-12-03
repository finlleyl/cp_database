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

var Module = fx.Options(

	user.Module,
	account.Module,

	strategy.Module,
	offer.Module,
	subscription.Module,
	trade.Module,

	statistics.Module,
	batchimport.Module,
	audit.Module,
)
