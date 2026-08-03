package repository

import "github.com/jmoiron/sqlx"

type ShareRepository struct {
	db *sqlx.DB
}

func NewShareRepository(db *sqlx.DB) *ShareRepository {
	return &ShareRepository{db: db}
}

func (r *ShareRepository) ListUsernames(boardID int) ([]string, error) {
	var usernames []string
	err := r.db.Select(&usernames,
		`SELECT u.username
		 FROM board_shares s
		 JOIN users u ON u.id = s.user_id
		 WHERE s.board_id = $1
		 ORDER BY u.username`,
		boardID,
	)
	if err != nil {
		return nil, err
	}
	return usernames, nil
}
