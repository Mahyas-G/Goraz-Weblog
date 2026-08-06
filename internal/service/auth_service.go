package service

import (
	"errors"
	"strings"

	"weblog/internal/auth"
	"weblog/internal/model"
	"weblog/internal/repository"
	"weblog/internal/validation"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

type AuthService struct {
	users    *repository.UserRepository
	sessions *auth.SessionStore
}

func NewAuthService(users *repository.UserRepository, sessions *auth.SessionStore) *AuthService {
	return &AuthService{users: users, sessions: sessions}
}

func (s *AuthService) Signup(username, password string) (*model.User, error) {
	username = strings.ToLower(strings.TrimSpace(username))

	if err := validation.ValidateUsername(username); err != nil {
		return nil, err
	}
	if err := validation.ValidatePassword(password); err != nil {
		return nil, err
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return s.users.Create(username, hash)
}

func (s *AuthService) Login(username, password string) (*model.Session, error) {
	username = strings.ToLower(strings.TrimSpace(username))

	user, err := s.users.FindByUsername(username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := auth.CheckPassword(user.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.sessions.Create(user.ID)
}

func (s *AuthService) Logout(token string) error {
	return s.sessions.Delete(token)
}

func (s *AuthService) CurrentUser(token string) (*model.User, error) {
	session, err := s.sessions.Load(token)
	if err != nil {
		return nil, err
	}
	return s.users.FindByID(session.UserID)
}
