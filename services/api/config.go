package main

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port uint32 `env:"PORT" envDefault:"4444"`
}

func NewConfig() (*Config, error) {
	config := Config{}

	err := env.Parse(&config)

	return &config, err
}
