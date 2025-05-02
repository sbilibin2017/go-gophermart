-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_balances (
    login VARCHAR(255) PRIMARY KEY,
    current DECIMAL(10, 2) NOT NULL DEFAULT 0,
    withdrawn BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (login) REFERENCES users(login) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_balances;
-- +goose StatementEnd
