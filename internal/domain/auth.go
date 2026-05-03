package domain

import (
	"errors"
	"time"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"

	AccessTokenTTL  = 24 * time.Hour
	RefreshTokenTTL = 7 * 24 * time.Hour
)

var ErrInvalidTokenType = errors.New("invalid token type")
