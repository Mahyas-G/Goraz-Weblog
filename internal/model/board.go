package model

import "time"

const (
	PrivacyPublic  = "public"
	PrivacyPrivate = "private"
)

type Board struct {
	ID             int       `db:"id"`
	Title          string    `db:"title"`
	Content        string    `db:"content"`
	ImagePath      *string   `db:"image_path"`
	AuthorID       int       `db:"author_id"`
	AuthorUsername string    `db:"author_username"`
	Privacy        string    `db:"privacy"`
	CreatedAt      time.Time `db:"created_at"`
}
