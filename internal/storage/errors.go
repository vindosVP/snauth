package storage

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrFailedToCreateUser = errors.New("failed to create user")
	ErrUserDoesNotExist   = errors.New("user does not exist")
	ErrFailedToGetUser    = errors.New("failed to get user")
)
