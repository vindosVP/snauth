package auth

import "github.com/pkg/errors"

var (
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
	ErrUserUnableToLogIn      = errors.New("user is unable to log in")
)
