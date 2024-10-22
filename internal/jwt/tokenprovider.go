package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"

	"github.com/vindosVP/snauth/internal/models"
)

type TokenProvider struct {
	serviceName string
	secret      []byte
	tokenTTL    time.Duration
	refreshTTL  time.Duration
}

func NewTokenProvider(secret []byte, tokenTTL, refreshTTL time.Duration) *TokenProvider {
	return &TokenProvider{secret: secret, tokenTTL: tokenTTL, refreshTTL: refreshTTL}
}

func (p *TokenProvider) ParseRefresh(refreshToken string) (int64, error) {
	f := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return p.secret, nil
	}
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, f)
	if err != nil {
		return 0, ErrInvalidToken
	}
	if !token.Valid {
		return 0, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, ErrInvalidToken
	}
	return claims.Id, nil
}

func (p *TokenProvider) NewPair(email string, id int64) (*models.TokenPair, error) {
	accessClaims := p.newAccessClaims(email, id)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(p.secret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign access token")
	}

	refreshClaims := p.newRefreshClaims(id)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(p.secret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign refresh token")
	}

	return &models.TokenPair{
		AccessToken:  accessString,
		RefreshToken: refreshString,
	}, nil
}

type Claims struct {
	jwt.RegisteredClaims
	Email string `json:"email,omitempty"`
	Id    int64  `json:"id"`
}

func (p *TokenProvider) newAccessClaims(email string, id int64) *Claims {
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    p.serviceName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Email: email,
		Id:    id,
	}
}

func (p *TokenProvider) newRefreshClaims(id int64) *Claims {
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    p.serviceName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Id: id,
	}
}
