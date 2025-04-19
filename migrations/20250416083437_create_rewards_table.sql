-- +goose Up
-- +goose StatementBegin
CREATE TABLE rewards (
    match VARCHAR(255) PRIMARY KEY,
    reward BIGINT NOT NULL,
    reward_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rewards;
-- +goose StatementEnd