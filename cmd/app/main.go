package main

import (
	_ "github.com/finlleyl/cp_database/docs"
	"github.com/finlleyl/cp_database/internal/config"
	"github.com/finlleyl/cp_database/internal/domain"
	"github.com/finlleyl/cp_database/internal/httpserver"
	"github.com/finlleyl/cp_database/internal/logger"
	"github.com/finlleyl/cp_database/internal/repository"
	"go.uber.org/fx"
)

// @title           Copy Trading API
// @version         1.0
// @description     API сервис для платформы копитрейдинга
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	fx.New(

		logger.Module,
		config.Module,
		repository.Module,

		domain.Module,

		httpserver.Module,
	).Run()
}
