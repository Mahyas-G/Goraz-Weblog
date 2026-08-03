package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"weblog/internal/model"
)

const sessionDuration = 7 * 24 * time.Hour

var ErrSessionNotFound = errors.New("session not found")

type SessionStore struct {
	db *sqlx.DB
}

func NewSessionStore(db *sqlx.DB) *SessionStore {
	return &SessionStore{db: db}
}

func (s *SessionStore) Create(userID int) (*model.Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	session := &model.Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionDuration),
	}

	_, err = s.db.Exec(
		`INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)`,
		session.Token, session.UserID, session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionStore) Load(token string) (*model.Session, error) {
	var session model.Session
	err := s.db.Get(&session,
		`SELECT token, user_id, expires_at FROM sessions WHERE token = $1`,
		token,
	)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.Delete(token)
		return nil, ErrSessionNotFound
	}

	return &session, nil
}

func (s *SessionStore) Delete(token string) error {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE token = $1`, token)
	return err
}

func (s *SessionStore) DeleteExpired() (int64, error) {
	result, err := s.db.Exec(`DELETE FROM sessions WHERE expires_at < now()`)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
