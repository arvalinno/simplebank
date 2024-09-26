package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrTokenInvalidClaims = errors.New("token has invalid claims")
	ErrInvalidToken       = errors.New("token is invalid")
	ErrExpiredToken       = errors.New("token is expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func (payload *Payload) Valid() error {
	expiredTime := payload.RegisteredClaims.ExpiresAt.Time

	if time.Now().After(expiredTime) {
		return ErrExpiredToken
	}
	return nil
}
