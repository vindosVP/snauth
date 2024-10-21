package jwt

import "errors"

var (
	ErrFailedToCreateTokens = errors.New("failed to create tokens")
	ErrInvalidToken         = errors.New("invalid token")
	ErrFailedToParseToken   = errors.New("failed to parse token")
)
