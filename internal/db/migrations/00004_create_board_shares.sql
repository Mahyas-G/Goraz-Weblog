-- +goose Up
CREATE TABLE board_shares (
    id       SERIAL PRIMARY KEY,
    board_id INTEGER NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    user_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(board_id, user_id)
);

-- +goose Down
DROP TABLE board_shares;
