package storage

import "github.com/pkg/errors"

var (
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrUserDoesNotExist  = errors.New("user does not exist")
)
