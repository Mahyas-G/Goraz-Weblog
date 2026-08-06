package validation

import (
	"errors"
	"regexp"

	"weblog/internal/model"
)

const (
	maxTitleLength    = 200
	maxContentLength  = 50000
	maxCommentLength  = 2000
	minUsernameLength = 3
	maxUsernameLength = 20
)

var usernamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

var ErrUsernameEmpty = errors.New("username is required")
var ErrUsernameTooShort = errors.New("username must be at least 3 characters")
var ErrUsernameTooLong = errors.New("username must be at most 20 characters")
var ErrInvalidUsername = errors.New("username must start with a letter and contain only English letters, numbers, and underscores")
var ErrPasswordTooShort = errors.New("password must be at least 8 characters")
var ErrTitleEmpty = errors.New("title is required")
var ErrTitleTooLong = errors.New("title must be at most 200 characters")
var ErrContentEmpty = errors.New("content is required")
var ErrContentTooLong = errors.New("content must be at most 50000 characters")
var ErrInvalidPrivacy = errors.New("privacy must be either public or private")
var ErrCommentEmpty = errors.New("comment is required")
var ErrCommentTooLong = errors.New("comment must be at most 2000 characters")

func ValidateUsername(username string) error {
	if username == "" {
		return ErrUsernameEmpty
	}
	if len(username) < minUsernameLength {
		return ErrUsernameTooShort
	}
	if len(username) > maxUsernameLength {
		return ErrUsernameTooLong
	}
	if !usernamePattern.MatchString(username) {
		return ErrInvalidUsername
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func ValidateTitle(title string) error {
	if title == "" {
		return ErrTitleEmpty
	}
	if len(title) > maxTitleLength {
		return ErrTitleTooLong
	}
	return nil
}

func ValidateContent(content string) error {
	if content == "" {
		return ErrContentEmpty
	}
	if len(content) > maxContentLength {
		return ErrContentTooLong
	}
	return nil
}

func ValidatePrivacy(privacy string) error {
	if privacy != model.PrivacyPublic && privacy != model.PrivacyPrivate {
		return ErrInvalidPrivacy
	}
	return nil
}

func ValidateComment(content string) error {
	if content == "" {
		return ErrCommentEmpty
	}
	if len(content) > maxCommentLength {
		return ErrCommentTooLong
	}
	return nil
}
