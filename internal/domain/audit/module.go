package audit

import (
	"go.uber.org/fx"
)

// Module provides all audit domain dependencies
var Module = fx.Options(
	fx.Provide(
		NewRepository,
		NewUseCase,
		NewHandler,
	),
)
