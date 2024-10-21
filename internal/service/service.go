package auth

import (
	"context"
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/vindosVP/snauth/internal/jwt"
	"github.com/vindosVP/snauth/internal/models"
	"github.com/vindosVP/snauth/internal/storage"
)

type UserStorage interface {
	CreateUser(ctx context.Context, email string, hPassword []byte, l *slog.Logger) (int64, error)
	UserByEmail(ctx context.Context, email string, l *slog.Logger) (*models.User, error)
	UserByID(ctx context.Context, id int64, l *slog.Logger) (*models.User, error)
	UserCanLogInByEmail(ctx context.Context, email string, l *slog.Logger) (bool, error)
	UserCanLogInById(ctx context.Context, id int64, l *slog.Logger) (bool, error)
}

type TokenProvider interface {
	NewPair(email string, id int64, l *slog.Logger) (*models.TokenPair, error)
	ParseRefresh(refreshToken string) (int64, error)
}

type Auth struct {
	us UserStorage
	t  TokenProvider
}

func New(us UserStorage, tp TokenProvider) *Auth {
	return &Auth{
		us: us,
		t:  tp,
	}
}

func (a *Auth) Register(ctx context.Context, email string, password string, l *slog.Logger) (int64, error) {
	hPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		l.Error("unable to hash password", slog.String("error", err.Error()))
		return 0, ErrFailedToRegisterUser
	}
	id, err := a.us.CreateUser(ctx, email, hPassword, l)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			l.Info("user already exists")
			return 0, ErrUserAlreadyExists
		}
		return 0, ErrFailedToRegisterUser
	}
	return id, nil
}

func (a *Auth) Login(ctx context.Context, email string, password string, l *slog.Logger) (*models.TokenPair, error) {
	u, err := a.us.UserByEmail(ctx, email, l)
	if err != nil {
		if errors.Is(err, storage.ErrUserDoesNotExist) {
			l.Info("user does not exist")
			return nil, ErrUserDoesNotExist
		}
		l.Error("unable to get user by email", slog.String("error", err.Error()))
		return nil, ErrFailedToLogIn
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.HPassword), []byte(password))
	if err != nil {
		l.Info("invalid login or password")
		return nil, ErrInvalidLoginOrPassword
	}
	if u.Banned || u.Deleted {
		l.Info("user is banned or deleted")
		return nil, ErrUserUnableToLogIn
	}
	tp, err := a.t.NewPair(u.Email, u.Id, l)
	if err != nil {
		l.Info("unable to create token pair", slog.String("error", err.Error()))
		return nil, ErrFailedToLogIn
	}
	return tp, nil
}

func (a *Auth) Refresh(ctx context.Context, refreshToken string, l *slog.Logger) (*models.TokenPair, error) {
	id, err := a.t.ParseRefresh(refreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			l.Info("invalid refresh token")
			return nil, ErrInvalidRefreshToken
		}
		l.Error("unable to parse refresh token", slog.String("error", err.Error()))
		return nil, ErrFailedToRefreshToken
	}
	u, err := a.us.UserByID(ctx, id, l)
	if err != nil {
		if errors.Is(err, storage.ErrUserDoesNotExist) {
			l.Info("user does not exist")
			return nil, ErrUserDoesNotExist
		}
		l.Error("unable to get user by id", slog.String("error", err.Error()))
		return nil, ErrFailedToRefreshToken
	}
	if u.Banned || u.Deleted {
		l.Info("user is banned or deleted")
		return nil, ErrUserUnableToLogIn
	}
	tp, err := a.t.NewPair(u.Email, u.Id, l)
	if err != nil {
		l.Info("unable to create token pair", slog.String("error", err.Error()))
		return nil, ErrFailedToRefreshToken
	}
	return tp, nil
}
