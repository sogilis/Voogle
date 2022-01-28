package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	RedisAddr string `env:"REDIS_ADDR,required"`
	RedisPwd  string `env:"REDIS_PWD,required"`
	RedisDB   uint32 `env:"REDIS_DB" envDefault:"0"`
}

func NewConfig() (Config, error) {
	config := Config{}

	err := env.Parse(&config)

	return config, err
}
