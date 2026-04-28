package validation

import "net/mail"

const (
	MaxEmailLength       = 320
	MinPasswordLength    = 8
	MaxPasswordLength    = 72
	MaxTitleLength       = 200
	MaxDescriptionLength = 4000
)

func IsValidEmail(s string) bool {
	if len(s) == 0 || len(s) > MaxEmailLength {
		return false
	}

	addr, err := mail.ParseAddress(s)
	if err != nil {
		return false
	}

	return addr.Address == s
}
