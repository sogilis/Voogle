package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port     uint32 `env:"PORT" envDefault:"4444"`
	UserAuth string `env:"USER_AUTH,required"`
	PwdAuth  string `env:"PWD_AUTH,required"`
	IsDev    bool   `env:"DEV_MODE" envDefault:"false"`
}

func NewConfig() (*Config, error) {
	config := Config{}

	err := env.Parse(&config)

	return &config, err
}
