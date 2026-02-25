package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test_jwt_secret"

func newTestService() *JWTService {
	return NewJWTService(testSecret)
}

func TestNewJWTService(t *testing.T) {
	t.Run("sets the key from secret", func(t *testing.T) {
		svc := NewJWTService(testSecret)
		if len(svc.key) == 0 {
			t.Fatal("expected key to be set")
		}
		if string(svc.key) != testSecret {
			t.Fatalf("expected key %q, got %q", testSecret, string(svc.key))
		}
	})
}

func TestJWTService_GenerateJWT(t *testing.T) {
	svc := newTestService()

	t.Run("returns a non-empty token string", func(t *testing.T) {
		tokenStr, err := svc.GenerateJWT(42)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tokenStr == "" {
			t.Fatal("expected a non-empty token string")
		}
	})

	t.Run("generated token contains correct user_id claim", func(t *testing.T) {
		userID := 99
		tokenStr, err := svc.GenerateJWT(userID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		claims, err := svc.ValidateJWT(tokenStr)
		if err != nil {
			t.Fatalf("expected valid token, got error: %v", err)
		}

		gotUserID := claims.UserID
		if gotUserID != userID {
			t.Fatalf("expected user_id %d, got %v", userID, gotUserID)
		}
	})

	t.Run("generated token contains a future expiration claim", func(t *testing.T) {
		tokenStr, err := svc.GenerateJWT(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		claims, err := svc.ValidateJWT(tokenStr)
		if err != nil {
			t.Fatalf("expected valid token, got error: %v", err)
		}

		exp := claims.ExpiresAt
		if !exp.After(time.Now()) {
			t.Fatal("expected expiration to be in the future")
		}
	})
}

func TestJWTService_ValidateJWT(t *testing.T) {
	svc := newTestService()

	t.Run("valid token returns claims without error", func(t *testing.T) {
		tokenStr, err := svc.GenerateJWT(7)
		if err != nil {
			t.Fatalf("expected no error generating token, got %v", err)
		}

		claims, err := svc.ValidateJWT(tokenStr)
		if err != nil {
			t.Fatalf("expected no error validating token, got %v", err)
		}
		if claims == nil {
			t.Fatal("expected non-nil claims")
		}
	})

	t.Run("tampered token returns error", func(t *testing.T) {
		tokenStr, err := svc.GenerateJWT(7)
		if err != nil {
			t.Fatalf("expected no error generating token, got %v", err)
		}

		_, err = svc.ValidateJWT(tokenStr + "tampered")
		if err == nil {
			t.Fatal("expected an error for tampered token, got nil")
		}
	})

	t.Run("token signed with different secret returns error", func(t *testing.T) {
		claims := jwt.MapClaims{
			"user_id": 1,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte("wrong_secret"))
		if err != nil {
			t.Fatalf("expected no error signing, got %v", err)
		}

		_, err = svc.ValidateJWT(tokenStr)
		if err == nil {
			t.Fatal("expected an error for token signed with wrong secret, got nil")
		}
	})

	t.Run("expired token returns error", func(t *testing.T) {
		claims := jwt.MapClaims{
			"user_id": 1,
			"exp":     time.Now().Add(-1 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(testSecret))
		if err != nil {
			t.Fatalf("expected no error signing, got %v", err)
		}

		_, err = svc.ValidateJWT(tokenStr)
		if err == nil {
			t.Fatal("expected an error for expired token, got nil")
		}
	})

	t.Run("empty token string returns error", func(t *testing.T) {
		_, err := svc.ValidateJWT("")
		if err == nil {
			t.Fatal("expected an error for empty token string, got nil")
		}
	})

	t.Run("two services with different secrets cannot validate each other's tokens", func(t *testing.T) {
		svc2 := NewJWTService("another_secret")
		tokenStr, err := svc2.GenerateJWT(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = svc.ValidateJWT(tokenStr)
		if err == nil {
			t.Fatal("expected error when validating token from a different service key")
		}
	})
}
