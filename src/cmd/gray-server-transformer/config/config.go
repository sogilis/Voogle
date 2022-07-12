package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port      uint32 `env:"PORT" envDefault:"50051"`
	LocalAddr string `env:"LOCAL_ADDR" envDefault:""`
	DevMode   bool   `env:"DEV_MODE" envDefault:"false"`

	S3Host    string `env:"S3_HOST" envDefault:""`
	S3AuthKey string `env:"S3_AUTH_KEY,required"`
	S3AuthPwd string `env:"S3_AUTH_PWD,required"`
	S3Bucket  string `env:"S3_BUCKET" envDefault:"voogle-video"`
	S3Region  string `env:"S3_REGION" envDefault:"eu-west-3"`

	ConsulHost string `env:"CONSUL_URL,required"`
	ConsulUser string `env:"CONSUL_USER" envDefault:""`
	ConsulPwd  string `env:"CONSUL_PWD" envDefault:""`
}

func NewConfig() (Config, error) {
	config := Config{}

	err := env.Parse(&config)

	return config, err
}
