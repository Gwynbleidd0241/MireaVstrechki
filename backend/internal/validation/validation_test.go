package validation

import (
	"strings"
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"valid simple", "user@example.com", true},
		{"valid with subdomain", "user@mail.example.com", true},
		{"valid with plus tag", "user+tag@example.com", true},
		{"valid with dots", "first.last@example.com", true},
		{"valid russian-looking", "petr.ivanov@company.ru", true},

		{"empty", "", false},
		{"no @", "userexample.com", false},
		{"no domain", "user@", false},
		{"no local part", "@example.com", false},
		{"only @", "@", false},
		{"with display name", "Petr <user@example.com>", false},
		{"trailing space", "user@example.com ", false},
		{"leading space", " user@example.com", false},
		{"contains space", "u ser@example.com", false},
		{"too long", strings.Repeat("a", 320) + "@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidEmail(tt.in)
			if got != tt.want {
				t.Errorf("IsValidEmail(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestLengthConstants(t *testing.T) {
	if MaxEmailLength != 320 {
		t.Errorf("MaxEmailLength = %d, want 320 (RFC 5321)", MaxEmailLength)
	}

	if MinPasswordLength != 8 {
		t.Errorf("MinPasswordLength = %d, want 8", MinPasswordLength)
	}

	if MaxPasswordLength != 72 {
		t.Errorf("MaxPasswordLength = %d, want 72 (bcrypt limit)", MaxPasswordLength)
	}

	if MaxTitleLength <= 0 || MaxDescriptionLength <= 0 {
		t.Error("title/description max lengths must be positive")
	}

	if MaxTitleLength >= MaxDescriptionLength {
		t.Error("title length should be less than description length")
	}
}
