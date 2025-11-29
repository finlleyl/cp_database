package batchimport

import (
	"go.uber.org/fx"
)

// Module provides all batch import domain dependencies
var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewUseCase,
		NewHandler,
	),
)

