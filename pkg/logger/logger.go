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
	var zl zerolog.Logger
	switch env {
	case envDev:
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		zl = zerolog.New(output)
		zl.Level(zerolog.DebugLevel)
	case envProd:
		zl = zerolog.New(os.Stdout)
		zl.Level(zerolog.InfoLevel)
	case envTest:
		zl = zerolog.New(os.Stdout)
		zl.Level(zerolog.Disabled)
	}
	return zl.With().Timestamp().Str("service", serviceName).Logger()
}
