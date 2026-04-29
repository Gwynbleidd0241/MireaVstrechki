package auth

import "testing"

func FuzzParseToken(f *testing.F) {
	f.Add("", "secret")
	f.Add("not-a-token", "secret")
	f.Add("a.b.c", "secret")
	f.Add("...", "secret")
	f.Add("\x00\x00\x00", "secret")

	if real, err := GenerateToken(42, "x@example.com", "admin", "secret"); err == nil {
		f.Add(real, "secret")
		f.Add(real, "wrong-secret")
		f.Add(real, "")
	}

	f.Fuzz(func(t *testing.T, token, secret string) {
		_, _ = ParseToken(token, secret)
	})
}

func FuzzGenerateToken(f *testing.F) {
	f.Add(int64(1), "user@example.com", "admin", "secret")
	f.Add(int64(0), "", "", "")
	f.Add(int64(-1), "\x00\x01\x02", "💀", "🔑")
	f.Add(int64(9223372036854775807), "x", "y", "z")

	f.Fuzz(func(t *testing.T, userID int64, email, role, secret string) {
		_, _ = GenerateToken(userID, email, role, secret)
	})
}
