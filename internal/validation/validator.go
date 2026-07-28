package validation

import (
	"errors"
	"strings"
)

var ErrUsernameEmpty = errors.New("username is required")
var ErrUsernameTooShort = errors.New("username must be at least 3 characters")
var ErrPasswordTooShort = errors.New("password must be at least 8 characters")

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return ErrUsernameEmpty
	}
	if len(username) < 3 {
		return ErrUsernameTooShort
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}
