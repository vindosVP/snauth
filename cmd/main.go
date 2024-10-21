package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/vindosVP/snauth/cmd/config"
	"github.com/vindosVP/snauth/internal/app"
	"github.com/vindosVP/snauth/pkg/logger/sl"
)

func main() {
	cfg := config.MustParse()
	l := sl.SetupLogger(cfg.Logger.ENV, cfg.ServiceName)

	a := app.New(l, cfg)
	go func() {
		a.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	a.GRPCServer.Stop()
	l.Info("Gracefully stopped")
}
