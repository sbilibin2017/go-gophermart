-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_withdraws (
    number VARCHAR(255) REFERENCES user_orders(number) ON DELETE CASCADE,
    sum BIGINT NOT NULL DEFAULT 0,
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_withdraws;
-- +goose StatementEnd

