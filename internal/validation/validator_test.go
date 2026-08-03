package validation

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateUsername(t *testing.T) {
	cases := []struct {
		name     string
		username string
		wantErr  error
	}{
		{"empty", "", ErrUsernameEmpty},
		{"whitespace only", "   ", ErrUsernameEmpty},
		{"too short", "ab", ErrUsernameTooShort},
		{"valid", "alice", nil},
		{"valid with surrounding whitespace", "  alice  ", nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateUsername(tc.username)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidateUsername(%q) = %v, want %v", tc.username, err, tc.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name     string
		password string
		wantErr  error
	}{
		{"too short", "short1", ErrPasswordTooShort},
		{"exactly minimum", "12345678", nil},
		{"long enough", "correcthorsebatterystaple", nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.password)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidatePassword(%q) = %v, want %v", tc.password, err, tc.wantErr)
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	cases := []struct {
		name    string
		title   string
		wantErr error
	}{
		{"empty", "", ErrTitleEmpty},
		{"valid", "My First Post", nil},
		{"too long", strings.Repeat("a", maxTitleLength+1), ErrTitleTooLong},
		{"exactly at limit", strings.Repeat("a", maxTitleLength), nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateTitle(tc.title)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidateTitle(...) = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestValidateContent(t *testing.T) {
	cases := []struct {
		name    string
		content string
		wantErr error
	}{
		{"empty", "", ErrContentEmpty},
		{"valid", "Some content", nil},
		{"too long", strings.Repeat("a", maxContentLength+1), ErrContentTooLong},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateContent(tc.content)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidateContent(...) = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestValidatePrivacy(t *testing.T) {
	cases := []struct {
		name    string
		privacy string
		wantErr error
	}{
		{"public", "public", nil},
		{"private", "private", nil},
		{"invalid", "everyone", ErrInvalidPrivacy},
		{"empty", "", ErrInvalidPrivacy},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePrivacy(tc.privacy)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidatePrivacy(%q) = %v, want %v", tc.privacy, err, tc.wantErr)
			}
		})
	}
}

func TestValidateComment(t *testing.T) {
	cases := []struct {
		name    string
		content string
		wantErr error
	}{
		{"empty", "", ErrCommentEmpty},
		{"valid", "nice post!", nil},
		{"too long", strings.Repeat("a", maxCommentLength+1), ErrCommentTooLong},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateComment(tc.content)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidateComment(...) = %v, want %v", err, tc.wantErr)
			}
		})
	}
}
