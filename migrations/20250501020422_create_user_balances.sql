-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_balances (
    login VARCHAR(255) REFERENCES users(login) ON DELETE CASCADE,
    current DECIMAL(20, 1) NOT NULL DEFAULT 0,
    withdrawn BIGINT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_balances;
-- +goose StatementEnd