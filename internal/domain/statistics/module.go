package statistics

import (
	"go.uber.org/fx"
)

// Module provides all statistics domain dependencies
var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewUseCase,
		NewHandler,
	),
)
