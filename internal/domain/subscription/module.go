package subscription

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewUseCase,
		NewHandler,
	),
)
