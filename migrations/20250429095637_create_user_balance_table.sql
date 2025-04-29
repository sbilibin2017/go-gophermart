-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_balances (
    login TEXT PRIMARY KEY REFERENCES users(login) ON DELETE CASCADE,
    current NUMERIC(20, 2) NOT NULL DEFAULT 0,
    withdrawn BIGINT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_balances;
-- +goose StatementEnd
