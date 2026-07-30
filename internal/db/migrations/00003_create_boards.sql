-- +goose Up
CREATE TABLE boards (
    id         SERIAL PRIMARY KEY,
    title      TEXT NOT NULL,
    content    TEXT NOT NULL,
    image_path TEXT,
    author_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    privacy    TEXT NOT NULL CHECK (privacy IN ('public', 'private')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE boards;
