package auth

import (
	"context"

	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/vindosVP/snauth/internal/jwt"
	"github.com/vindosVP/snauth/internal/models"
	"github.com/vindosVP/snauth/internal/storage"
)

type UserStorage interface {
	CreateUser(ctx context.Context, email string, hPassword []byte) (int64, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
	UserByID(ctx context.Context, id int64) (*models.User, error)
}

type TokenProvider interface {
	NewPair(email string, id int64) (*models.TokenPair, error)
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

func (a *Auth) Register(ctx context.Context, email string, password string) (int64, error) {
	hPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.Wrap(err, "failed to hash password")
	}
	id, err := a.us.CreateUser(ctx, email, hPassword)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			return 0, ErrUserAlreadyExists
		}
		return 0, errors.Wrap(err, "failed create user")
	}
	return id, nil
}

func (a *Auth) Login(ctx context.Context, email string, password string) (*models.TokenPair, error) {
	u, err := a.us.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserDoesNotExist) {
			return nil, ErrInvalidLoginOrPassword
		}
		return nil, errors.Wrap(err, "failed to get user by email")
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.HPassword), []byte(password))
	if err != nil {
		return nil, ErrInvalidLoginOrPassword
	}
	if u.Banned || u.Deleted {
		return nil, ErrUserUnableToLogIn
	}
	tp, err := a.t.NewPair(u.Email, u.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create token pair")
	}
	return tp, nil
}

func (a *Auth) Refresh(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	id, err := a.t.ParseRefresh(refreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, errors.Wrap(err, "failed to parse refresh token")
	}
	u, err := a.us.UserByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserDoesNotExist) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, errors.Wrap(err, "failed to get user by id")
	}
	if u.Banned || u.Deleted {
		return nil, ErrUserUnableToLogIn
	}
	tp, err := a.t.NewPair(u.Email, u.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create token pair")
	}
	return tp, nil
}
