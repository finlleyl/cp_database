package account

import (
	"go.uber.org/fx"
)

// Module provides all account domain dependencies
var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewUseCase,
		NewHandler,
	),
)

