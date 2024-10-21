// Package sl initializes and configures slog logger
package sl

import (
	"log/slog"
	"os"

	"github.com/vindosVP/snauth/pkg/logger/handlers/slogdiscard"
	"github.com/vindosVP/snauth/pkg/logger/handlers/slogpretty"
)

const (
	// envLocal is a local environment
	envLocal = "local"

	// envLocal is a development environment
	envDev = "dev"

	// envLocal is a production environment
	envProd = "prod"

	// envTest is a test environment
	envTest = "test"
)

// SetupLogger configures the logger depending on environment
func SetupLogger(env string, serviceName string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		opts := slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}

		handler := opts.NewPrettyHandler(os.Stdout)

		log = slog.New(handler)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slogdiscard.NewDiscardHandler(),
		)
	}
	return log.With(slog.String("service", serviceName))
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
