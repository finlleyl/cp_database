package user

import (
	"go.uber.org/fx"
)

// Module provides all user domain dependencies
var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewUseCase,
		NewHandler,
	),
)

