package storage

import (
	"context"

	"github.com/pkg/errors"

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

func (us *UserStorage) CreateUser(ctx context.Context, email string, hPassword []byte) (int64, error) {
	u, err := us.s.UserByEmail(ctx, email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}
	if u != nil {
		return 0, ErrUserAlreadyExists
	}
	id, err := us.s.CreateUser(ctx, email, hPassword)
	if err != nil {
		return 0, errors.Wrap(err, "failed to save user")
	}
	return id, nil
}

func (us *UserStorage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	u, err := us.s.UserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserDoesNotExist
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user by email in db")
	}
	return u, nil
}

func (us *UserStorage) UserByID(ctx context.Context, id int64) (*models.User, error) {
	u, err := us.s.UserByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserDoesNotExist
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user by id in db")
	}
	return u, nil
}
