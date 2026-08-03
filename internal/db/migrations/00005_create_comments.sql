-- +goose Up
CREATE TABLE comments (
    id         SERIAL PRIMARY KEY,
    board_id   INTEGER NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    author_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_comments_board_id ON comments(board_id);

-- +goose Down
DROP TABLE comments;
