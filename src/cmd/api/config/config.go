package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port  uint32 `env:"PORT" envDefault:"4444"`
	IsDev bool   `env:"DEV_MODE" envDefault:"false"`

	UserAuth string `env:"USER_AUTH,required"`
	PwdAuth  string `env:"PWD_AUTH,required"`

	S3Host    string `env:"S3_HOST" envDefault:""`
	S3AuthKey string `env:"S3_AUTH_KEY,required"`
	S3AuthPwd string `env:"S3_AUTH_PWD,required"`
	S3Bucket  string `env:"S3_BUCKET" envDefault:"voogle-video"`
	S3Region  string `env:"S3_REGION" envDefault:"eu-west-3"`

	RedisAddr string `env:"REDIS_ADDR,required"`
	RedisPwd  string `env:"REDIS_PWD,required"`
}

func NewConfig() (Config, error) {
	config := Config{}

	err := env.Parse(&config)

	return config, err
}
