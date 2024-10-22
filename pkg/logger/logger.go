package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	envDev  = "dev"
	envProd = "prod"
	envTest = "test"
)

func SetupLogger(env string, serviceName string) zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	zl := zerolog.New(output)
	switch env {
	case envDev:
		zl.Level(zerolog.DebugLevel)
	case envProd:
		zl.Level(zerolog.InfoLevel)
	case envTest:
		zl.Level(zerolog.Disabled)
	}
	return zl.With().Timestamp().Str("service", serviceName).Str("env", env).Logger()
}
