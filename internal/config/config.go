package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
)

type Config struct {
	HTTPPort string `env:"HTTP_PORT" default:"8080"`
	POSTGRES_DB string `env:"POSTGRES_DB" default:"postgres"`
	POSTGRES_USER string `env:"POSTGRES_USER" default:"postgres"`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD" default:"postgres"`
}

func loadConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

var Module = fx.Provide(loadConfig)