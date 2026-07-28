package model

import "time"

type Session struct {
	Token     string    `db:"token"`
	UserID    int       `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
}
