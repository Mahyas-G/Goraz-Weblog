package model

import "time"

type Comment struct {
	ID             int       `db:"id"`
	BoardID        int       `db:"board_id"`
	AuthorID       int       `db:"author_id"`
	AuthorUsername string    `db:"author_username"`
	Content        string    `db:"content"`
	CreatedAt      time.Time `db:"created_at"`
}
