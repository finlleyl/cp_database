package main

import (
	"github.com/finlleyl/cp_database/internal/config"
	"github.com/finlleyl/cp_database/internal/domain"
	"github.com/finlleyl/cp_database/internal/httpserver"
	"github.com/finlleyl/cp_database/internal/logger"
	"github.com/finlleyl/cp_database/internal/repository"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Infrastructure modules
		logger.Module,
		config.Module,
		repository.Module,

		// Domain modules (all domains aggregated)
		domain.Module,

		// HTTP server module
		httpserver.Module,
	).Run()
}
