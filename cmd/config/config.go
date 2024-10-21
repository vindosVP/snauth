package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	DB          DB
	Token       Token
	GRPC        GRPC
	Logger      Logger
	ServiceName string `env:"SERVICE_NAME" envDefault:"auth"`
}

type DB struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	Database string `env:"DB_DATABASE"`
}

type Token struct {
	Secret     string        `env:"TOKEN_SECRET"`
	TokenTTL   time.Duration `env:"TOKEN_TTL"`
	RefreshTTL time.Duration `env:"REFRESH_TTL"`
}

type GRPC struct {
	Port    int           `env:"GRPC_PORT"`
	Timeout time.Duration `env:"GRPC_TIMEOUT"`
}

type Logger struct {
	ENV string `env:"LOG_ENV" envDefault:"dev"`
}

func MustParse() *Config {
	cfg := &Config{}
	err := env.Parse(cfg, env.Options{RequiredIfNoDef: true})
	if err != nil {
		panic(fmt.Errorf("filed to parse config: %w", err))
	}
	return cfg
}
