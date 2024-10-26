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
	SetDeletedToUser(ctx context.Context, userId int64, isDeleted bool) (bool, error)
	SetBannedToUser(ctx context.Context, userId int64, isBanned bool) (bool, error)
	SetAdminToUser(ctx context.Context, userId int64, isAdmin bool) (bool, error)
}

type TokenProvider interface {
	NewPair(email string, id int64, isAdmin bool) (*models.TokenPair, error)
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

func (a *Auth) SetDeleted(ctx context.Context, id int64, deleted bool) (bool, error) {
	_, err := a.us.UserByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserDoesNotExist) {
			return false, ErrUserDoesNotExist
		}
		return false, errors.Wrap(err, "failed to get user by id")
	}
	isDeleted, err := a.us.SetDeletedToUser(ctx, id, deleted)
	if err != nil {
		return false, errors.Wrap(err, "failed to set deleted to user")
	}
	return isDeleted, nil
}

func (a *Auth) SetBanned(ctx context.Context, id int64, banned bool) (bool, error) {
	_, err := a.us.UserByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserDoesNotExist) {
			return false, ErrUserDoesNotExist
		}
		return false, errors.Wrap(err, "failed to get user by id")
	}
	isBanned, err := a.us.SetBannedToUser(ctx, id, banned)
	if err != nil {
		return false, errors.Wrap(err, "failed to set banned to user")
	}
	return isBanned, nil
}

func (a *Auth) SetAdmin(ctx context.Context, id int64, admin bool) (bool, error) {
	_, err := a.us.UserByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserDoesNotExist) {
			return false, ErrUserDoesNotExist
		}
		return false, errors.Wrap(err, "failed to get user by id")
	}
	isAdmin, err := a.us.SetAdminToUser(ctx, id, admin)
	if err != nil {
		return false, errors.Wrap(err, "failed to set admin to user")
	}
	return isAdmin, nil
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
	if u.IsBanned || u.IsDeleted {
		return nil, ErrUserUnableToLogIn
	}
	tp, err := a.t.NewPair(u.Email, u.Id, u.IsAdmin)
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
	if u.IsBanned || u.IsDeleted {
		return nil, ErrUserUnableToLogIn
	}
	tp, err := a.t.NewPair(u.Email, u.Id, u.IsAdmin)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create token pair")
	}
	return tp, nil
}
