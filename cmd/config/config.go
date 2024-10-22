package config

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Config struct {
	DB          DB     `json:"db"`
	Token       Token  `json:"token"`
	GRPC        GRPC   `json:"gRPC"`
	Logger      Logger `json:"logger"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"auth" json:"serviceName"`
}

type DB struct {
	Host     string `env:"DB_HOST" json:"host"`
	Port     int    `env:"DB_PORT" json:"port"`
	Username string `env:"DB_USERNAME" json:"-"`
	Password string `env:"DB_PASSWORD" json:"-"`
	Database string `env:"DB_DATABASE" json:"database"`
}

type Token struct {
	Secret     string        `env:"TOKEN_SECRET" json:"-"`
	TokenTTL   time.Duration `env:"TOKEN_TTL" json:"token_ttl"`
	RefreshTTL time.Duration `env:"REFRESH_TTL" json:"refresh_ttl"`
}

type GRPC struct {
	Port    int           `env:"GRPC_PORT" json:"port"`
	Timeout time.Duration `env:"GRPC_TIMEOUT" json:"timeout"`
}

type Logger struct {
	ENV string `env:"LOG_ENV" envDefault:"dev" json:"env"`
}

func MustParse() *Config {
	cfg := &Config{}
	err := env.Parse(cfg, env.Options{RequiredIfNoDef: true})
	if err != nil {
		panic(errors.Wrap(err, "filed to parse config"))
	}
	return cfg
}
