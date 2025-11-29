package logger

import (
    "go.uber.org/fx"
    "go.uber.org/zap"
    "go.uber.org/fx/fxevent"
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
