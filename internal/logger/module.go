package logger

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	return cfg.Build()
}

var Module = fx.Options(
	fx.Provide(NewLogger),

	fx.WithLogger(func(l *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: l}
	}),
)
