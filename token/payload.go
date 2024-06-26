package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (payload *Payload) Validator() error {
	expiredTime := payload.RegisteredClaims.ExpiresAt.Time

	if time.Now().After(expiredTime) {
		return ErrExpiredToken
	}
	return nil
}
