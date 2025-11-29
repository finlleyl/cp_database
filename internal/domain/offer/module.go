package offer

import (
	"go.uber.org/fx"
)

// Module provides all offer domain dependencies
var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewUseCase,
		NewHandler,
	),
)

