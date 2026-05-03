package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kerhael/accounting/internal/domain"
)

type JWTService struct {
	key []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{key: []byte(secret)}
}

func (s *JWTService) GenerateAccessToken(userID int) (string, error) {
	return s.generateToken(userID, domain.AccessTokenType, domain.AccessTokenTTL)
}

func (s *JWTService) GenerateRefreshToken(userID int) (string, error) {
	return s.generateToken(userID, domain.RefreshTokenType, domain.RefreshTokenTTL)
}

func (s *JWTService) GenerateTokenPair(userID int) (accessToken string, refreshToken string, err error) {
	accessToken, err = s.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *JWTService) generateToken(userID int, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.key)
}

func (s *JWTService) ValidateJWT(tokenStr string) (*CustomClaims, error) {
	return s.validateToken(tokenStr, domain.AccessTokenType)
}

func (s *JWTService) ValidateRefreshToken(tokenStr string) (*CustomClaims, error) {
	return s.validateToken(tokenStr, domain.RefreshTokenType)
}

func (s *JWTService) validateToken(tokenStr string, expectedTokenType string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return s.key, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if claims.TokenType != expectedTokenType {
		return nil, domain.ErrInvalidTokenType
	}

	return claims, nil
}
