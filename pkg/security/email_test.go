package security

import (
	"testing"
)

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercases uppercase letters",
			input:    "User@Example.COM",
			expected: "user@example.com",
		},
		{
			name:     "trims leading and trailing spaces",
			input:    "  user@example.com  ",
			expected: "user@example.com",
		},
		{
			name:     "lowercases and trims simultaneously",
			input:    "  USER@EXAMPLE.COM  ",
			expected: "user@example.com",
		},
		{
			name:     "already normalized email is unchanged",
			input:    "user@example.com",
			expected: "user@example.com",
		},
		{
			name:     "empty string stays empty",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeEmail(tt.input)
			if got != tt.expected {
				t.Errorf("NormalizeEmail(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	t.Run("valid email returns nil", func(t *testing.T) {
		err := ValidateEmail("user@example.com")
		if err != nil {
			t.Errorf("expected nil for valid email, got %v", err)
		}
	})

	t.Run("valid email with uppercase letters returns nil (normalized)", func(t *testing.T) {
		err := ValidateEmail("USER@EXAMPLE.COM")
		if err != nil {
			t.Errorf("expected nil for uppercased valid email, got %v", err)
		}
	})

	t.Run("valid email with surrounding spaces returns nil (normalized)", func(t *testing.T) {
		err := ValidateEmail("  user@example.com  ")
		if err != nil {
			t.Errorf("expected nil for email with spaces, got %v", err)
		}
	})

	t.Run("empty string returns error", func(t *testing.T) {
		err := ValidateEmail("")
		if err == nil {
			t.Fatal("expected an error for empty email, got nil")
		}
	})

	t.Run("whitespace-only string returns error", func(t *testing.T) {
		err := ValidateEmail("   ")
		if err == nil {
			t.Fatal("expected an error for whitespace-only email, got nil")
		}
	})

	t.Run("missing @ returns error", func(t *testing.T) {
		err := ValidateEmail("userexample.com")
		if err == nil {
			t.Fatal("expected an error for email without @, got nil")
		}
	})

	t.Run("missing domain returns error", func(t *testing.T) {
		err := ValidateEmail("user@")
		if err == nil {
			t.Fatal("expected an error for email without domain, got nil")
		}
	})

	t.Run("plain text returns error", func(t *testing.T) {
		err := ValidateEmail("not-an-email")
		if err == nil {
			t.Fatal("expected an error for plain text, got nil")
		}
	})
}
