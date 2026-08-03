package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"weblog/internal/model"
)

var ErrBoardNotFound = errors.New("board not found")

type BoardRepository struct {
	db *sqlx.DB
}

func NewBoardRepository(db *sqlx.DB) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) Create(title, content string, imagePath *string, authorID int, privacy string, sharedUserIDs []int) (*model.Board, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var board model.Board
	err = tx.Get(&board,
		`INSERT INTO boards (title, content, image_path, author_id, privacy)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, title, content, image_path, author_id, privacy, created_at`,
		title, content, imagePath, authorID, privacy,
	)
	if err != nil {
		return nil, err
	}

	for _, userID := range sharedUserIDs {
		if _, err := tx.Exec(
			`INSERT INTO board_shares (board_id, user_id) VALUES ($1, $2)`,
			board.ID, userID,
		); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &board, nil
}

func (r *BoardRepository) ListVisibleToUser(userID int) ([]model.Board, error) {
	var boards []model.Board
	err := r.db.Select(&boards,
		`SELECT b.id, b.title, b.content, b.image_path, b.author_id,
		        u.username AS author_username, b.privacy, b.created_at
		 FROM boards b
		 JOIN users u ON u.id = b.author_id
		 WHERE b.privacy = 'public'
		    OR b.author_id = $1
		    OR b.id IN (SELECT board_id FROM board_shares WHERE user_id = $1)
		 ORDER BY b.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func (r *BoardRepository) FindVisibleToUser(id, userID int) (*model.Board, error) {
	var board model.Board
	err := r.db.Get(&board,
		`SELECT b.id, b.title, b.content, b.image_path, b.author_id,
		        u.username AS author_username, b.privacy, b.created_at
		 FROM boards b
		 JOIN users u ON u.id = b.author_id
		 WHERE b.id = $1
		   AND (b.privacy = 'public'
		        OR b.author_id = $2
		        OR b.id IN (SELECT board_id FROM board_shares WHERE user_id = $2))`,
		id, userID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBoardNotFound
	}
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *BoardRepository) DeleteOwned(id, authorID int) (bool, error) {
	result, err := r.db.Exec(`DELETE FROM boards WHERE id = $1 AND author_id = $2`, id, authorID)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
