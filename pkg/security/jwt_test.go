package security

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test_jwt_secret"

func setupJWT() {
	InitJWT(testSecret)
}

func TestInitJWT(t *testing.T) {
	t.Run("sets the JWT key", func(t *testing.T) {
		InitJWT(testSecret)
		if len(jwtKey) == 0 {
			t.Fatal("expected jwtKey to be set after InitJWT")
		}
		if string(jwtKey) != testSecret {
			t.Fatalf("expected jwtKey to be %q, got %q", testSecret, string(jwtKey))
		}
	})
}

func TestGenerateJWT(t *testing.T) {
	setupJWT()

	t.Run("returns a non-empty token string", func(t *testing.T) {
		tokenStr, err := GenerateJWT(42)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tokenStr == "" {
			t.Fatal("expected a non-empty token string")
		}
	})

	t.Run("generated token contains correct user_id claim", func(t *testing.T) {
		userID := 99
		tokenStr, err := GenerateJWT(userID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		claims, err := ValidateJWT(tokenStr)
		if err != nil {
			t.Fatalf("expected valid token, got error: %v", err)
		}

		gotUserID, ok := claims["user_id"]
		if !ok {
			t.Fatal("expected 'user_id' claim to be present")
		}
		// JWT numbers are float64 when decoded via MapClaims
		if int(gotUserID.(float64)) != userID {
			t.Fatalf("expected user_id %d, got %v", userID, gotUserID)
		}
	})

	t.Run("generated token contains an expiration claim", func(t *testing.T) {
		tokenStr, err := GenerateJWT(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		claims, err := ValidateJWT(tokenStr)
		if err != nil {
			t.Fatalf("expected valid token, got error: %v", err)
		}

		exp, ok := claims["exp"]
		if !ok {
			t.Fatal("expected 'exp' claim to be present")
		}
		expTime := time.Unix(int64(exp.(float64)), 0)
		if !expTime.After(time.Now()) {
			t.Fatal("expected expiration to be in the future")
		}
	})
}

func TestValidateJWT(t *testing.T) {
	setupJWT()

	t.Run("valid token returns claims without error", func(t *testing.T) {
		tokenStr, err := GenerateJWT(7)
		if err != nil {
			t.Fatalf("expected no error generating token, got %v", err)
		}

		claims, err := ValidateJWT(tokenStr)
		if err != nil {
			t.Fatalf("expected no error validating token, got %v", err)
		}
		if claims == nil {
			t.Fatal("expected non-nil claims")
		}
	})

	t.Run("tampered token returns error", func(t *testing.T) {
		tokenStr, err := GenerateJWT(7)
		if err != nil {
			t.Fatalf("expected no error generating token, got %v", err)
		}

		tampered := tokenStr + "tampered"
		_, err = ValidateJWT(tampered)
		if err == nil {
			t.Fatal("expected an error for tampered token, got nil")
		}
	})

	t.Run("token signed with different secret returns error", func(t *testing.T) {
		// Sign with a different key
		claims := jwt.MapClaims{
			"user_id": 1,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte("wrong_secret"))
		if err != nil {
			t.Fatalf("expected no error signing, got %v", err)
		}

		_, err = ValidateJWT(tokenStr)
		if err == nil {
			t.Fatal("expected an error for token signed with wrong secret, got nil")
		}
	})

	t.Run("expired token returns error", func(t *testing.T) {
		claims := jwt.MapClaims{
			"user_id": 1,
			"exp":     time.Now().Add(-1 * time.Hour).Unix(), // expired 1 hour ago
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(testSecret))
		if err != nil {
			t.Fatalf("expected no error signing, got %v", err)
		}

		_, err = ValidateJWT(tokenStr)
		if err == nil {
			t.Fatal("expected an error for expired token, got nil")
		}
	})

	t.Run("empty token string returns error", func(t *testing.T) {
		_, err := ValidateJWT("")
		if err == nil {
			t.Fatal("expected an error for empty token string, got nil")
		}
	})
}
