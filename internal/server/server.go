package server

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authv1 "github.com/vindosVP/snauth/gen/go"
	"github.com/vindosVP/snauth/internal/models"
	auth "github.com/vindosVP/snauth/internal/service"
)

type Auth interface {
	Register(ctx context.Context, email string, password string, l *slog.Logger) (int64, error)
	Login(ctx context.Context, email string, password string, l *slog.Logger) (*models.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string, l *slog.Logger) (*models.TokenPair, error)
}

type server struct {
	authv1.UnimplementedAuthServer
	auth Auth
	l    *slog.Logger
}

func Register(gRPCServer *grpc.Server, auth Auth, l *slog.Logger) {
	authv1.RegisterAuthServer(gRPCServer, &server{auth: auth, l: l})
}

func requestID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("no metadata")
	}
	MDreqId := md.Get("requestID")
	if len(MDreqId) == 0 {
		return "", errors.New("no request id")
	}
	return MDreqId[0], nil
}

func (s *server) Register(ctx context.Context, in *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	reqId, err := requestID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	l := s.l.With(slog.String("requestID", reqId))
	id, err := s.auth.Register(ctx, in.GetEmail(), in.GetPassword(), l)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		l.Error("Failed to register user", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authv1.RegisterResponse{UserId: id}, nil
}

func (s *server) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	reqId, err := requestID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	l := s.l.With(slog.String("requestID", reqId))
	tokenPair, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), l)
	if err != nil {
		if errors.Is(err, auth.ErrUserDoesNotExist) || errors.Is(err, auth.ErrInvalidLoginOrPassword) ||
			errors.Is(err, auth.ErrUserUnableToLogIn) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		l.Error("Failed to register user", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authv1.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s *server) Refresh(ctx context.Context, in *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	reqId, err := requestID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	l := s.l.With(slog.String("requestID", reqId))
	tokenPair, err := s.auth.Refresh(ctx, in.GetRefreshToken(), l)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRefreshToken) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, auth.ErrUserDoesNotExist) || errors.Is(err, auth.ErrUserUnableToLogIn) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		l.Error("Failed to refresh token", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authv1.RefreshResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
