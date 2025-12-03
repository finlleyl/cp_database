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

		logger.Module,
		config.Module,
		repository.Module,

		domain.Module,

		httpserver.Module,
	).Run()
}
