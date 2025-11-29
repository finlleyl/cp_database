package trade

import (
	"go.uber.org/fx"
)

// Module provides all trade domain dependencies
var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewCopiedTradeRepository,
		NewUseCase,
		NewHandler,
	),
)

