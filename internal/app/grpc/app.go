package grpc

import (
	"context"
	"fmt"
	"net"
	"slices"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vindosVP/snauth/internal/server"
)

type App struct {
	l          zerolog.Logger
	gRPCServer *grpc.Server
	port       int
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	a.l.Info().Str("addr", l.Addr().String()).Msg("grpc server started")
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func (a *App) Stop() {
	a.l.Info().Msg("stopping grpc server")
	a.gRPCServer.GracefulStop()
}

func New(log zerolog.Logger, a server.Auth, port int) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
		logging.WithDisableLoggingFields(
			logging.ComponentFieldKey, logging.ServiceFieldKey, logging.SystemTag[0],
		),
	}
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error().Any("panic", p).Msg("Recovered from panic")
			return status.Errorf(codes.Internal, "internal error")
		}),
	}
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))
	server.Register(gRPCServer, a, log)
	return &App{
		l:          log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func InterceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		fields = clearSensitiveData(fields)
		l := l.With().Fields(fields).Logger()
		switch lvl {
		case logging.LevelDebug:
			l.Debug().Msg(msg)
		case logging.LevelInfo:
			l.Info().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func clearSensitiveData(fields []any) []any {
	id := slices.Index(fields, "grpc.request.content")
	if id != -1 {
		fields = slices.Delete(fields, id+1, id+2)
	}
	id = slices.Index(fields, "grpc.response.content")
	if id != -1 {
		fields = slices.Delete(fields, id+1, id+2)
	}
	return fields
}
