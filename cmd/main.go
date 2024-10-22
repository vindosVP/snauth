package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/vindosVP/snauth/cmd/config"
	"github.com/vindosVP/snauth/internal/app"
	"github.com/vindosVP/snauth/pkg/logger"
)

func main() {
	cfg := config.MustParse()
	l := logger.SetupLogger(cfg.Logger.ENV, cfg.ServiceName)
	l.Info().Interface("config", cfg).Msg("configuration loaded")

	a := app.New(l, cfg)
	go func() {
		a.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	a.GRPCServer.Stop()
	l.Info().Msg("gracefully stopped")
}
