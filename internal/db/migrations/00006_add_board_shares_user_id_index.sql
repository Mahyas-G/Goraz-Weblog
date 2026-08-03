-- +goose Up
CREATE INDEX idx_board_shares_user_id ON board_shares(user_id);

-- +goose Down
DROP INDEX idx_board_shares_user_id;
