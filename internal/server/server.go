package server

import (
	"context"

	"github.com/pkg/errors"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authv1 "github.com/vindosVP/snauth/gen/go"
	"github.com/vindosVP/snauth/internal/models"
	auth "github.com/vindosVP/snauth/internal/service"
)

type Auth interface {
	Register(ctx context.Context, email string, password string) (int64, error)
	Login(ctx context.Context, email string, password string) (*models.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*models.TokenPair, error)
}

type server struct {
	authv1.UnimplementedAuthServer
	auth Auth
	l    zerolog.Logger
}

func Register(gRPCServer *grpc.Server, auth Auth, l zerolog.Logger) {
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
		s.l.Error().Err(err).Msg("failed to extract request ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	l := s.l.With().Str("requestID", reqId).Str("email", in.GetEmail()).Logger()
	l.Info().Msg("registering user")
	id, err := s.auth.Register(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			l.Info().Msg("user already exists")
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		l.Error().Stack().Err(err).Msg("failed to register user")
		return nil, status.Error(codes.Internal, "failed to register user")
	}
	l.Info().Msg("registered user successfully")
	return &authv1.RegisterResponse{UserId: id}, nil
}

func (s *server) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	reqId, err := requestID(ctx)
	if err != nil {
		s.l.Error().Err(err).Msg("failed to extract request ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	l := s.l.With().Str("requestID", reqId).Str("email", in.GetEmail()).Logger()
	l.Info().Msg("logging user in")
	tokenPair, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidLoginOrPassword) {
			l.Info().Msg("invalid login or password")
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, auth.ErrUserUnableToLogIn) {
			l.Info().Msg("user is deleted or banned")
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		l.Error().Stack().Err(err).Msg("failed to log in user")
		return nil, status.Error(codes.Internal, "failed to log in user")
	}
	l.Info().Msg("logged in user successfully")
	return &authv1.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s *server) Refresh(ctx context.Context, in *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	reqId, err := requestID(ctx)
	if err != nil {
		s.l.Error().Err(err).Msg("failed to extract request ID")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	l := s.l.With().Str("requestID", reqId).Logger()
	l.Info().Msg("refreshing token")
	tokenPair, err := s.auth.Refresh(ctx, in.GetRefreshToken())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRefreshToken) {
			l.Info().Msg("invalid refresh token")
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, auth.ErrUserUnableToLogIn) {
			l.Info().Msg("unable to refresh token")
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		l.Error().Stack().Err(err).Msg("failed to refresh token")
		return nil, status.Error(codes.Internal, "failed to refresh token")
	}
	return &authv1.RefreshResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
