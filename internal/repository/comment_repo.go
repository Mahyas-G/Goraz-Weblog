package repository

import (
	"github.com/jmoiron/sqlx"

	"weblog/internal/model"
)

type CommentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(boardID, authorID int, content string) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Get(&comment,
		`INSERT INTO comments (board_id, author_id, content)
		 VALUES ($1, $2, $3)
		 RETURNING id, board_id, author_id, content, created_at`,
		boardID, authorID, content,
	)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) ListByBoard(boardID int) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Select(&comments,
		`SELECT c.id, c.board_id, c.author_id, u.username AS author_username, c.content, c.created_at
		 FROM comments c
		 JOIN users u ON u.id = c.author_id
		 WHERE c.board_id = $1
		 ORDER BY c.created_at ASC`,
		boardID,
	)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
