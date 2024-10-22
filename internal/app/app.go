package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/vindosVP/snauth/cmd/config"
	"github.com/vindosVP/snauth/internal/app/grpc"
	"github.com/vindosVP/snauth/internal/jwt"
	auth "github.com/vindosVP/snauth/internal/service"
	"github.com/vindosVP/snauth/internal/storage"
	"github.com/vindosVP/snauth/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpc.App
}

func New(log zerolog.Logger, cfg *config.Config) *App {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, postgresConn(cfg))
	if err != nil {
		panic(fmt.Errorf("could not connect to postgres: %w", err))
	}
	p := postgres.New(pool)
	us := storage.NewUserStorage(p)
	tp := jwt.NewTokenProvider([]byte(cfg.Token.Secret), cfg.Token.TokenTTL, cfg.Token.RefreshTTL)
	authService := auth.New(us, tp)
	grpcApp := grpc.New(log, authService, cfg.GRPC.Port)
	return &App{
		GRPCServer: grpcApp,
	}
}

func postgresConn(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
	)
}
