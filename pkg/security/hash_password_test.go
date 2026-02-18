package security

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	t.Run("returns a non-empty hash", func(t *testing.T) {
		hash, err := HashPassword("mysecretpassword")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if hash == "" {
			t.Fatal("expected a non-empty hash")
		}
	})

	t.Run("hash differs from original password", func(t *testing.T) {
		password := "mysecretpassword"
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if hash == password {
			t.Fatal("hash should not equal the original password")
		}
	})

	t.Run("same password produces different hashes (bcrypt salting)", func(t *testing.T) {
		password := "mysecretpassword"
		hash1, err := HashPassword(password)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		hash2, err := HashPassword(password)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if hash1 == hash2 {
			t.Fatal("two hashes of the same password should differ due to bcrypt salting")
		}
	})
}

func TestCheckPassword(t *testing.T) {
	t.Run("correct password returns nil error", func(t *testing.T) {
		password := "correctpassword"
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("expected no error hashing, got %v", err)
		}

		err = CheckPassword(password, hash)
		if err != nil {
			t.Fatalf("expected nil error for correct password, got %v", err)
		}
	})

	t.Run("wrong password returns an error", func(t *testing.T) {
		password := "correctpassword"
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("expected no error hashing, got %v", err)
		}

		err = CheckPassword("wrongpassword", hash)
		if err == nil {
			t.Fatal("expected an error for wrong password, got nil")
		}
	})

	t.Run("empty password against valid hash returns error", func(t *testing.T) {
		hash, err := HashPassword("somepassword")
		if err != nil {
			t.Fatalf("expected no error hashing, got %v", err)
		}

		err = CheckPassword("", hash)
		if err == nil {
			t.Fatal("expected an error for empty password, got nil")
		}
	})
}
