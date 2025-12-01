package subscription

import (
	"go.uber.org/fx"
)

// Module provides all subscription domain dependencies
var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewUseCase,
		NewHandler,
	),
)
