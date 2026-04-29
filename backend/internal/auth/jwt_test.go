package auth

import (
	"strings"
	"testing"
)

func TestGenerateAndParseToken_RoundTrip(t *testing.T) {
	secret := "test-secret"

	token, err := GenerateToken(42, "user@example.com", "admin", secret)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	if !strings.Contains(token, ".") {
		t.Errorf("expected JWT format with dots, got %q", token)
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("UserID = %d, want 42", claims.UserID)
	}

	if claims.Email != "user@example.com" {
		t.Errorf("Email = %q, want user@example.com", claims.Email)
	}

	if claims.Role != "admin" {
		t.Errorf("Role = %q, want admin", claims.Role)
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken(1, "x@example.com", "admin", "secret-a")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	if _, err := ParseToken(token, "secret-b"); err == nil {
		t.Error("expected error when parsing with wrong secret")
	}
}

func TestParseToken_Tampered(t *testing.T) {
	token, err := GenerateToken(1, "x@example.com", "employee", "secret")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	// Tamper with the signature segment.
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("unexpected token format: %s", token)
	}

	tampered := parts[0] + "." + parts[1] + ".tampered"

	if _, err := ParseToken(tampered, "secret"); err == nil {
		t.Error("expected error when parsing tampered token")
	}
}

func TestParseToken_Garbage(t *testing.T) {
	cases := []string{
		"",
		"not-a-token",
		"abc.def",
		"...",
	}

	for _, c := range cases {
		if _, err := ParseToken(c, "secret"); err == nil {
			t.Errorf("expected error parsing %q", c)
		}
	}
}

func TestGenerateToken_HasExpiration(t *testing.T) {
	token, err := GenerateToken(1, "x@example.com", "admin", "secret")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	claims, err := ParseToken(token, "secret")
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}

	if claims.ExpiresAt == nil {
		t.Fatal("ExpiresAt is nil, expected token to have expiration")
	}

	if claims.IssuedAt == nil {
		t.Fatal("IssuedAt is nil")
	}

	if !claims.ExpiresAt.After(claims.IssuedAt.Time) {
		t.Error("ExpiresAt should be after IssuedAt")
	}
}
