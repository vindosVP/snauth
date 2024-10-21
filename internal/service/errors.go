package auth

import "errors"

var (
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrFailedToRegisterUser   = errors.New("failed to create user")
	ErrUserDoesNotExist       = errors.New("user does not exist")
	ErrFailedToLogIn          = errors.New("failed to log in user")
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
	ErrFailedToRefreshToken   = errors.New("failed to refresh token")
	ErrUserUnableToLogIn      = errors.New("user is unable to log in")
)
