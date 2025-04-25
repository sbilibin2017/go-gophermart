-- +goose Up
-- +goose StatementBegin
CREATE TABLE accrual_reward_mechanic (
    match VARCHAR(100) PRIMARY KEY,
    reward BIGINT NOT NULL,
    reward_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accrual_reward_mechanic;
-- +goose StatementEnd
