package validation

import (
	"net/mail"
	"strings"
	"testing"
)

func FuzzIsValidEmail(f *testing.F) {
	f.Add("user@example.com")
	f.Add("first.last+tag@mail.example.com")
	f.Add("")
	f.Add("a")
	f.Add("@")
	f.Add("user@@example.com")
	f.Add("user @example.com")
	f.Add("user@example.com\n")
	f.Add("\x00user@example.com")
	f.Add("Petr <user@example.com>")
	f.Add(strings.Repeat("a", 400) + "@example.com")

	f.Fuzz(func(t *testing.T, s string) {
		valid := IsValidEmail(s)

		if !valid {
			return
		}

		addr, err := mail.ParseAddress(s)
		if err != nil {
			t.Errorf("IsValidEmail(%q) = true, but mail.ParseAddress failed: %v", s, err)
			return
		}

		if addr.Address != s {
			t.Errorf("IsValidEmail(%q) = true, but parsed addr.Address = %q (display-name leaked)",
				s, addr.Address)
		}

		if len(s) > MaxEmailLength {
			t.Errorf("IsValidEmail(%q) = true, but len > MaxEmailLength=%d", s, MaxEmailLength)
		}
	})
}
