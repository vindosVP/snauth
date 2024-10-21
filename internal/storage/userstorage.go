package storage

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/vindosVP/snauth/internal/models"
)

type Storage interface {
	CreateUser(ctx context.Context, email string, hPassword []byte) (int64, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
	UserByID(ctx context.Context, id int64) (*models.User, error)
}

type UserStorage struct {
	s Storage
}

func NewUserStorage(s Storage) *UserStorage {
	return &UserStorage{s}
}

func (us *UserStorage) CreateUser(ctx context.Context, email string, hPassword []byte, l *slog.Logger) (int64, error) {
	u, err := us.s.UserByEmail(ctx, email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}
	if u != nil {
		return 0, ErrUserAlreadyExists
	}
	id, err := us.s.CreateUser(ctx, email, hPassword)
	if err != nil {
		l.Error("failed to create user", slog.String("error", err.Error()))
		return 0, ErrFailedToCreateUser
	}
	return id, nil
}

func (us *UserStorage) UserByEmail(ctx context.Context, email string, l *slog.Logger) (*models.User, error) {
	u, err := us.s.UserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserDoesNotExist
	}
	if err != nil {
		l.Error("failed to get user by email", slog.String("error", err.Error()))
		return nil, ErrFailedToGetUser
	}
	return u, nil
}

func (us *UserStorage) UserByID(ctx context.Context, id int64, l *slog.Logger) (*models.User, error) {
	u, err := us.s.UserByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserDoesNotExist
	}
	if err != nil {
		l.Error("failed to get user by id", slog.String("error", err.Error()))
		return nil, ErrFailedToGetUser
	}
	return u, nil
}

func (us *UserStorage) UserCanLogInByEmail(ctx context.Context, email string, l *slog.Logger) (bool, error) {
	u, err := us.s.UserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		l.Error("failed to get user by email", slog.String("error", err.Error()))
		return false, ErrFailedToGetUser
	}
	return !u.Banned && !u.Deleted, nil
}

func (us *UserStorage) UserCanLogInById(ctx context.Context, id int64, l *slog.Logger) (bool, error) {
	u, err := us.s.UserByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		l.Error("failed to get user by id", slog.String("error", err.Error()))
		return false, ErrFailedToGetUser
	}
	return !u.Banned && !u.Deleted, nil
}
