package strategy

import (
	"go.uber.org/fx"
)

// Module provides all strategy domain dependencies
var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewUseCase,
		NewHandler,
	),
)

